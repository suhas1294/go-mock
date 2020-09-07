package main

import (
	"github.com/suhas1294/go-mock/controllers"
	"net/http"
)

func main() {
	http.HandleFunc("/create_mock", controllers.MockController.CreateMock)
	http.HandleFunc("/get_mock/", controllers.MockController.GetMock)
	http.HandleFunc("/mock/", controllers.MockController.Mock)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

/**

err := json.Unmarshal([]byte(jsonData), &goData)
	if err != nil {
		fmt.Println("Some error while unmarshallnig\n\n")
		panic(err)
	}
	fmt.Println(goData.JsonPayload)
	// fmt.Println(string(goData.JsonPayload)) // json.RawMessage to string conversion not possible

	bp, err := json.MarshalIndent(goData.JsonPayload, "", "  ")
	fmt.Println(string(bp))
*/
