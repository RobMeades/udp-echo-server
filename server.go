/* Entry point for UDP echo server.
 *
 * This code based on the TCP echo server also in this repo.
 */
 
package main

import (
    "net"
    "fmt"
    "os"
    "flag"
    "log"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

var numPackets int

// Command-line flags
var pPort = flag.String ("p", "", "the UDP port number to listen on.")
var Usage = func() {
    fmt.Fprintf(os.Stderr, "\n%s: run the UDP echo server.  Usage:\n", os.Args[0])
        flag.PrintDefaults()
    }

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Entry point
func main() {
    var readErr error
    var numBytesIn int
    var pRemoteUdpAddr *net.UDPAddr
    line := make([]byte, 1024)

    // Deal with the command-line parameters
    flag.Parse()
    
    // Set up logging
    log.SetFlags(log.LstdFlags)
    
    if *pPort != "" {
        // Say what we're doing
        fmt.Printf("Echoing UDP packets received on port %s.\n", *pPort)
        
        // Set up the server
        pLocalUdpAddr, err := net.ResolveUDPAddr("udp", ":" + *pPort)
        if (err == nil) && (pLocalUdpAddr != nil) {
            pServer, err := net.ListenUDP("udp", pLocalUdpAddr)
            if err == nil {
                if err == nil {                    
                    for numBytesIn, pRemoteUdpAddr, readErr = pServer.ReadFromUDP(line); (readErr == nil) && (numBytesIn > 0); numBytesIn, pRemoteUdpAddr, readErr = pServer.ReadFromUDP(line) {
                        numBytesOut, writeErr := pServer.WriteToUDP(line[:numBytesIn], pRemoteUdpAddr)
                        log.Printf("%d: %v <-> %v: %s\n", numPackets, pLocalUdpAddr, pRemoteUdpAddr, line[:numBytesIn])
                        if writeErr != nil {
                            log.Printf("Error, only %d of %d byte(s) could be echoed (%s).\n", numBytesOut, numBytesIn, writeErr.Error())    
                        }
                    }    
                    if readErr != nil {
                        fmt.Printf("Error reading from port %v (%s).\n", pLocalUdpAddr, readErr.Error())
                    } else {
                        log.Printf("UDP read on port %v returned when it should not.\n", pLocalUdpAddr)    
                    }
                } else {
                    fmt.Printf("Couldn't disable read deadline (%s).\n", err.Error())
                }
            } else {
                fmt.Printf("Couldn't start UDP server on port %s (%s).\n", *pPort, err.Error())
            }            
        } else {
            fmt.Printf("'%s' is not a valid UDP address (%s).\n", *pPort, err.Error())
        }            
    } else {
        fmt.Printf("Must specify a port number.\n")
        flag.PrintDefaults()
        os.Exit(-1)
    }
}

// End Of File
