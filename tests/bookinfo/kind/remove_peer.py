#!/usr/bin/env python3
import os,time
import subprocess as sp
import sys
import argparse

proj_dir = os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname( os.path.abspath(__file__)))))
sys.path.insert(0,f'{proj_dir}')

from tests.utils.mbgAux import runcmd, runcmdb, printHeader, getPodName
from tests.utils.kind.kindAux import useKindCluster

def removePeer(mbgName, mbgctlName, peerName):
    useKindCluster(mbgName)
    mbgctlPod = getPodName("mbgctl")
    printHeader(f"\n\Remove {peerName} peer to {mbgName}")
    runcmd(f'kubectl exec -i {mbgctlPod} -- ./mbgctl remove peer --myid {mbgctlName} --id {peerName}')

############################### MAIN ##########################
if __name__ == "__main__":
    #parameters 
    mbgName     = "mbg1"
    mbgCtlName  = "mbgctl1"
    peerName    = "mbg3"
    print(f'Working directory {proj_dir}')
    os.chdir(proj_dir)

    removePeer(mbgName,mbgCtlName,peerName)
    

    