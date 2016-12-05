/* Entry point for UDP echo server.
 *
 * This code based:
 * https://gist.github.com/paulsmith/775764
 */
 
package main

import (
    "net"
    "fmt"
    "os"
    "flag"
    "bufio"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// File handle
var pFile *os.File
var numPackets int

// Command-line flags
var pPort = flag.String ("p", "", "the port number to listen on.")
var pFileName = flag.String ("f", "", "the file name to write the receive requests to.")
var Usage = func() {
    fmt.Fprintf(os.Stderr, "\n%s: run the UDP echo server.  Usage:\n", os.Args[0])
        flag.PrintDefaults()
    }

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Entry point
func main() {
    var err error
    var line []byte

    // Deal with the command-line parameters
    flag.Parse()
    
    if *pPort != "" {
        // Open the output file for append
        if *pFileName != "" {
            pFile, err = os.OpenFile(*pFileName, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0666)
        }    
        
        if err == nil {        
            // Say what we're doing
            fmt.Printf("Echoing UDP packets received on port %s", *pPort)
            if (pFile != nil) {
                fmt.Printf(" and writing received packets to \"%s\"", pFile.Name())
            }
            fmt.Printf(".\n")
            
            // Set up the server
            pLocalUdpAddr, err := net.ResolveUDPAddr("udp", ":" + *pPort)
            if err == nil {
                pServer, err := net.ListenUDP("udp", pLocalUdpAddr)
                if err == nil {
                    for numBytesIn, pRemoteUdpAddr, readErr := pServer.ReadFromUDP(line); readErr == nil; {
                        if numBytesIn > 0 {
                            fmt.Printf("%d: %v <-> %v.\n", numPackets, pLocalUdpAddr, pRemoteUdpAddr)                    
                            numBytesOut, writeErr := pServer.WriteToUDP(line, pRemoteUdpAddr)
                            if writeErr != nil {
                                fmt.Printf("Error, only %d of %d byte(s) could be echoed (%s).\n", numBytesOut, numBytesIn, writeErr.Error())    
                            }
                            if pFile != nil {
                                writer := bufio.NewWriter(pFile)
                                fmt.Fprintf(writer, "%d: %v <-> %v \"%s\".\n", numPackets, pLocalUdpAddr, pRemoteUdpAddr, line)
                            }
                        }
                    }    
                    fmt.Printf("Unable to read from port %v (%s).\n", pLocalUdpAddr)    
                } else {                
                    fmt.Printf("Couldn't start UDP server on port %s (%s).\n", *pPort, err.Error())
                }            
            } else {
                fmt.Printf("'%s' is not a valid UDP address (%s).\n", *pPort, err.Error())
            }            
        } else {
            fmt.Printf("Couldn't open file %s (%s).\n", *pFileName, err.Error())
            os.Exit(-1)
        }
    } else {
        fmt.Printf("Must specify a port number.\n")
        flag.PrintDefaults()
        os.Exit(-1)
    }
}

// End Of File
