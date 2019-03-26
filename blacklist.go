package main

import (
        "encoding/json"
        "fmt"

        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
        "errors"
        "strconv"
        "bytes"
        "encoding/pem"
        "crypto/x509"
)

type Blacklist struct{}

// ===================================================================================
// Main
// ===================================================================================
func main() {
        err := shim.Start(new(Blacklist))
        if err != nil {
                fmt.Printf("Error starting Simple chaincode: %s", err)
        }
}

func (fc *Blacklist) Init(stub shim.ChaincodeStubInterface) pb.Response {
        return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (fc *Blacklist) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        function, args := stub.GetFunctionAndParameters()
        fmt.Println("invoke is running " + function)

        // Handle different functions
        if function == "uploadBlacklist" {
                return fc.uploadBlacklist(stub, args)
        } else if function == "readBlacklist" {
                return fc.readBlacklist(stub, args)
        } 

        fmt.Println("invoke did not find func: " + function) //error
        return shim.Error("Received unknown function invocation")
}

func (fc *Blacklist) write(stub shim.ChaincodeStubInterface, key string,value string) error {
        var err error
        if len(key) <= 0 {
                return errors.New("1st argument(key) must be a non-empty string")
        }

        if len(value) <= 0 {
                return errors.New("2nd argument(value) must be a non-empty string")
        }

        var blob []byte = []byte(value)
        fmt.Println("write key:"+key+",value:"+value)
        err = stub.PutState(key, blob)
        return err
}

func (fc *Blacklist) uploadBlacklist(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2, key value optype ")
        }

        var key = args[0]
        var value = args[1]
        fmt.Println("uploadBlacklist:key:"+key+"value:"+value)
        if len(key) <= 0 {
                return shim.Error("1st argument(key) must be a non-empty string")
        }

        if len(value) <= 0 {
                return shim.Error("2nd argument(value) must be a non-empty string")
        }

        valAsbytes, err := fc.read(stub, key)
        if(valAsbytes != nil){
            return shim.Success(nil)
        }
        ......
        return shim.Success(nil)
}

func (fc *Blacklist) readBlacklist(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        var err error
        var jsonResp string
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1, key  optype ")
        }
        key := args[0]
        valAsbytes, err := fc.read(stub,key)
        if err != nil {
                jsonResp = err.Error()
                return shim.Error(jsonResp)
        } else if valAsbytes == nil {
                return shim.Success(valAsbytes)
        }else{
            result := string(valAsbytes[:])
            var data = make(map[string]string)
            err = json.Unmarshal(valAsbytes,&data)
            if(err != nil){
                jsonResp = "{\"Error\":\"unmarshal result fail,key: "+args[0] +" ,result:"+ result + "\"}"
                fmt.Println("unmarshal result error")
                fmt.Println(err)
                return shim.Error(jsonResp)
            }
            ......
        }
}

func (fc *Blacklist) history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        var jsonResp string

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. expect 1, key")
        }

        var key = args[0]

        if len(key) <= 0 {
                return shim.Error("1st argument(key) must be a non-empty string")
        }

        var itor shim.HistoryQueryIteratorInterface
        itor, err = stub.GetHistoryForKey(key)
        if err != nil {
                jsonResp = "{\"Error\":\"Failed to get history data for " + key + "\"}"
                return shim.Error(jsonResp)
        } else if itor == nil {
                jsonResp = "{\"Error\":\" history itor is null\"}"
                return shim.Error(jsonResp)
        }

        var keyArray []interface{} = make([]interface{}, 0)
        for itor.HasNext() {
                kmdf, err := itor.Next()
                if err != nil {
                        jsonResp = "{\"Error\":\" traverse history itor error" + err.Error() + " \"}"
                        return shim.Error(jsonResp)
                }
                keyArray = append(keyArray, kmdf)
        }

        jsonBolb,err := json.Marshal(keyArray)
        if err != nil {
                jsonResp = "{\"Error\":\" json marshal error" + err.Error() + " \"}"
                return shim.Error(jsonResp)
        }

        return shim.Success(jsonBolb)

}

func (blacklist *Blacklist) getUserName(stub shim.ChaincodeStubInterface) string{
   creatorByte,_:= stub.GetCreator()
   certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
   if certStart == -1 {
      fmt.Errorf("No certificate found")
   }
   certText := creatorByte[certStart:]
   bl, _ := pem.Decode(certText)
   if bl == nil {
      fmt.Errorf("Could not decode the PEM structure")
   }

   cert, err := x509.ParseCertificate(bl.Bytes)
   if err != nil {
      fmt.Errorf("ParseCertificate failed")
   }
   uname:=cert.Subject.CommonName
   fmt.Println("Name:"+uname)
   return uname
}
