/**********************************************************/
/* mTLS Forwader : This is created per service-pair connections.
/**********************************************************/
// Generate Certificates
// openssl req -newkey rsa:2048   -new -nodes -x509   -days 3650   -out ~/mtls/tcnode7_cert.pem   -keyout ~/mtls/tcnode7_key.pem   -subj "/C=US/ST=California/L=mbg/O=ibm/OU=dev/CN=tcnode7" -addext "subjectAltName = IP:10.20.20.2"
// openssl req -newkey rsa:2048   -new -nodes -x509   -days 3650   -out ~/mtls/tcnode6_cert.pem   -keyout ~/mtls/tcnode6_key.pem   -subj "/C=US/ST=California/L=mbg/O=ibm/OU=dev/CN=tcnode6" -addext "subjectAltName = IP:10.20.20.1"

// Workflow of mTLS forwarder usage
// After Expose of a service at MBG 1 run the following APIs :
//    1) StartLocalService for the exported service at other remote application (for e.g. App 2)
//    2) When LocalService receives an accepted connection from APP 2, Do an Connect API to APP 1
//    3) MBG1 starts a StartReceiverService with the necessary details such as endpoint, and sends it as Connect Response

package mbgDataplane

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"strings"
	"time"

	"github.com/go-chi/chi"
)

type MbgMtlsForwarder struct {
	Name           string
	Connection     net.Conn
	mtlsConnection net.Conn
	ChiRouter      *chi.Mux
}

type connDialer struct {
	c net.Conn
}

func (cd connDialer) Dial(network, addr string) (net.Conn, error) {
	return cd.c, nil
}

//var clog = logrus.WithField("component", "mbgDataplane/mTLSForwarder")

// Start mtls Forwarder on a specific mtls target
// targetIPPort in the format of <ip:port>
// connect is set to true on a client side
func (m *MbgMtlsForwarder) StartMtlsForwarderClient(targetIPPort, name, rootCA, certificate, key, ServerName string, endpointConn net.Conn) {
	clog.Infof("Starting to initialize mTLS Forwarder for MBG Dataplane at /mbgData/%s", m.Name)

	connectMbg := "https://" + targetIPPort + "/mbgData/" + name

	m.Connection = endpointConn
	m.Name = name

	//Create TCP connection with TLS handshake
	TLSClientConfig := m.CreateTlsConfig(rootCA, certificate, key, ServerName)
	mtls_conn, err := tls.Dial("tcp", targetIPPort, TLSClientConfig)
	if err != nil {
		clog.Infof("Error in connecting.. %+v", err)
		m.CloseConnection()
		return
	}

	//clog.Debugln("mTLS Debug Check:", m.certDebg(targetIPPort, name, tlsConfig))

	TlsConnectClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: TLSClientConfig,
			DialTLS:         connDialer{mtls_conn}.Dial,
		},
	}
	req, err := http.NewRequest(http.MethodGet, connectMbg, nil)
	if err != nil {
		clog.Infof("Failed to create new request %v", err)
		m.CloseConnection()
		return
	}
	resp, err := TlsConnectClient.Do(req)
	if err != nil {
		clog.Infof("Error in Tls Connection %v", err)
		m.CloseConnection()
		return
	}

	m.mtlsConnection = mtls_conn
	clog.Infof("mtlS Connection Established Resp:%s(%d) to Target: %s", resp.Status, resp.StatusCode, connectMbg)

	//From forwarder to other MBG
	go m.mtlsDispatch()
	//From source to forwarder
	go m.dispatch()
	clog.Infof("Starting mTLS Forwarder client for MBG Dataplane at /mbgData/%s  to target %s with certs(%s,%s)", m.Name, targetIPPort, certificate, key)

}

// Start mtls Forwarder server on a specific mtls target
// Register handling function (for hijack the connection) and start dispatch to destination
func (m *MbgMtlsForwarder) StartMtlsForwarderServer(targetIPPort, name, rootCA, certificate, key string, endpointConn net.Conn) {
	clog.Infof("Starting to initialize mTLS Forwarder for MBG Dataplane at /mbgData/%s", m.Name)

	// Register function for handling the dataplane traffic
	clog.Infof("Register new handle func to address =%s", "/mbgData/"+name)
	m.ChiRouter.Get("/mbgData/"+name, m.mbgConnectHandler)

	connectMbg := "https://" + targetIPPort + "/mbgData/" + name

	clog.Infof("Connect MBG Target =%s", connectMbg)
	m.Connection = endpointConn
	m.Name = name

	go m.dispatch()
	clog.Infof("Starting mTLS Forwarder server for MBG Dataplane at /mbgData/%s  to target %s with certs(%s,%s)", m.Name, targetIPPort, certificate, key)

}

func (m *MbgMtlsForwarder) mbgConnectHandler(w http.ResponseWriter, r *http.Request) {
	clog.Infof("Received mbgConnect (%s) from %s", m.Name, r.RemoteAddr)

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "server doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	//Hijack the connection
	conn, _, err := hj.Hijack()
	if err != nil {
		clog.Infof("Hijacking failed %v", err)
		return
	}
	conn.Write([]byte{})
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n")

	clog.Infof("Connection Hijacked  %v->%v", conn.RemoteAddr().String(), conn.LocalAddr().String())

	m.mtlsConnection = conn
	clog.Infof("Starting to dispatch mtls Connection")
	go m.mtlsDispatch()
}

func (m *MbgMtlsForwarder) mtlsDispatch() error {
	bufData := make([]byte, maxDataBufferSize)
	var err error
	for {
		numBytes, err := m.mtlsConnection.Read(bufData)
		if err != nil {
			if err == io.EOF {
				err = nil //Ignore EOF error
			} else {
				clog.Infof("mtlsDispatch: Read error %v\n", err)
			}
			break
		}
		m.Connection.Write(bufData[:numBytes])
	}
	clog.Infof("Initiating end of mtls connection(%s)", m.Name)
	m.CloseConnection()
	if err == io.EOF {
		return nil
	} else {
		return err
	}
}

func (m *MbgMtlsForwarder) dispatch() error {
	bufData := make([]byte, maxDataBufferSize)
	var err error
	for {
		numBytes, err := m.Connection.Read(bufData)
		if err != nil {
			if err == io.EOF {
				err = nil //Ignore EOF error
			} else {
				clog.Errorf("Dispatch: Read error %v  connection: (local:%s Remote:%s)->,(local: %s Remote%s) ", err,
					m.Connection.LocalAddr(), m.Connection.RemoteAddr(), m.mtlsConnection.LocalAddr(), m.mtlsConnection.RemoteAddr())

			}
			break
		}
		if m.mtlsConnection == nil {
			clog.Info("Start Waiting for mtlsConnection") //start infinite loop
			for m.mtlsConnection == nil {
				time.Sleep(time.Microsecond)
			}
			clog.Info("Finish Waiting for mtlsConnection ") //Finish infinite loop
		}
		_, err = m.mtlsConnection.Write(bufData[:numBytes])
		if err != nil {
			clog.Errorf("Dispatch: Write error %v  connection: (local:%s Remote:%s)->,(local: %s Remote%s) ", err,
				m.Connection.LocalAddr(), m.Connection.RemoteAddr(), m.mtlsConnection.LocalAddr(), m.mtlsConnection.RemoteAddr())
			break
		}
	}
	clog.Infof("Initiating end of connection(%s)", m.Name)
	m.CloseConnection()
	if err == io.EOF {
		return nil
	} else {
		return err
	}
}

func (m *MbgMtlsForwarder) CloseConnection() {
	if m.Connection != nil {
		m.Connection.Close()
	}
	if m.mtlsConnection != nil {
		m.mtlsConnection.Close()
	}

}

func CloseMtlsServer(ip string) {
	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr: ip,
	}
	server.Shutdown(context.Background())
}

// Get rootCA, certificate, key  and create tls config
func (m *MbgMtlsForwarder) CreateTlsConfig(rootCA, certificate, key, ServerName string) *tls.Config {
	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		clog.Fatalf("LoadX509KeyPair -%v \ncertificate: %v \nkey:%v", err, certificate, key)
	}

	// Create a CA certificate pool and add ca to it
	caCert, err := ioutil.ReadFile(rootCA)
	if err != nil {
		clog.Fatalf("ReadFile certificate %v :%v", rootCA, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	TLSConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		ServerName:   ServerName,
	}
	return TLSConfig
}

// method for debug only -use to debug mtls connection
func (m *MbgMtlsForwarder) certDebg(target, name string, tlsConfig *tls.Config) string {
	clog.Infof("Starting tls debug to addr %v name %v", target, name)
	conn, err := tls.Dial("tcp", target, tlsConfig)
	if err != nil {
		panic("Server doesn't support SSL certificate err: " + err.Error())
	}
	ip := strings.Split(target, ":")[0]
	err = conn.VerifyHostname(ip)
	if err != nil {
		panic("Hostname doesn't match with certificate: " + err.Error())
	}
	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
	clog.Infof("Issuer: %s\nExpiry: %v\n", conn.ConnectionState().PeerCertificates[0].Issuer, expiry.Format(time.RFC850))
	conn.Close()
	return "Debug succeed"
}
