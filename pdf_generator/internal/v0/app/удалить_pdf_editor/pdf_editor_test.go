package delete_pdf_editor

import (
	"fmt"
	"github.com/manyakrus/starter/common/v0/micro"
	"github.com/manyakrus/starter/pdf_generator/internal/v0/app/programdir"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"
)

func test(j int) {
	ProgramDir := programdir.ProgramDir()
	filename := ProgramDir + "test.html"

	map1 := make(map[string]interface{})
	map1["name"] = "Никитин А.В."

	dir1 := filepath.Dir(filename)

	start := time.Now()

	for i := 1; i <= 10; i++ {
		//println(time.Now().String())

		filename2 := dir1 + micro.SeparatorFile() + "test_ready" + strconv.Itoa(j) + ".pdf"

		err := CreatePDF1(filename, filename2, map1)
		if err != nil {
			log.Error("удалить_pdf_editor.TestCreatePDF1() error: ", err)
		}
	}

	end := time.Since(start)
	fmt.Println(j, end.String())

	wg.Done()
}

var wg = &sync.WaitGroup{}

func TestCreatePDF1(t *testing.T) {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		j := i
		go test(j)
	}
	wg.Wait()

}
