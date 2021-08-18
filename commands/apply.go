package commands

import (
	"bytes"
	"fmt"
	"github.com/discless/discless-cli/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var ApplyCmd = &cobra.Command{
	Use: "apply [bot] [function configuration]",
	Short: "Create a new function with given name",
	Args: cobra.MinimumNArgs(2),
	RunE: FApply,
}

func FApply(c *cobra.Command, args []string) error {
	file, err := ioutil.ReadFile(args[1])
	if err != nil {
		return err
	}
	config := &types.Config{}
	yaml.Unmarshal(file, config)

	for function, ap := range config.Functions {
		PostApply(function,ap, args[0])
	}

	return nil
}

func PostApply(name string, function types.Function, bot string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("Name",name)
	writer.WriteField("Function", function.Function)
	writer.WriteField("Category", function.Category)
	writer.WriteField("Bot",bot)

	fw, err := writer.CreateFormFile("Function",function.File)
	if err != nil {
		return err
	}
	file, err := os.Open(function.File)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", "http://localhost:6969/function", bytes.NewReader(body.Bytes()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, _ := client.Do(req)

	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}

	fmt.Println("Succesfully uploaded the", name, "command")

	return nil
}
