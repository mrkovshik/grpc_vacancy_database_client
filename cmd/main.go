package main

import (
	"context"
	proto "github.com/mrkovshik/grpc_vacancy_database/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	_ "github.com/lib/pq"
)

var mainMenu string = "\n*****************************\n - Если хотите посмотреть всю таблицу вакансий, наберите \"посмотреть\", \n - Если хотите найти вакансию по названию наберите \"найти\"\n - Если хотите добавить строку - наберите \"добавить\", \n - Если хотите удалить вакансию из базы, наберите \"удалить\", \n - Если хотите выйти из программы, наберите \"выход\"\n*****************************\n"
func main() {

	err := mainDialog()
	if err != nil {
		fmt.Println(err)
	}
}

func gprcConnect ()(proto.DBServerClient,context.Context,error){
	ctx := context.TODO()
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("dial error")
		return nil,ctx, err
	}
	// defer conn.Close()
	client := proto.NewDBServerClient(conn)
	return client,ctx, err

}
func deleteDialog () (int,bool) {
	var vacId int
	var err error
	fmt.Println("Введите идентификационный номер вакансии, которую хотите удалить, либо наберите \"назад\" для выхода в предыдущее меню")
	scanner := bufio.NewScanner(os.Stdin)
	OuterLoop:
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			if scanner.Text() == "назад" {
				return 0, false
			}
			searchKey:= scanner.Text()
			vacId, err = strconv.Atoi(searchKey)
			if err != nil {
				fmt.Println("Ошибка ввода данных (идентификационный номер должен быть числом), повторите ввод")
				continue OuterLoop
			}

	fmt.Printf("Вы уверены, что хотите удалить из базы вакансию с номером %v? Наберите \"да\" или \"нет\"", vacId)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			if scanner.Text() == "да" {
				break
			}
			if scanner.Text() == "нет" {
				fmt.Println("Введите идентификационный номер вакансии, которую хотите удалить, либо наберите \"назад\" для выхода в предыдущее меню")
				continue OuterLoop
			}
			fmt.Println("Наберите \"да\" или \"нет\"")
			}
		}
	}

			return vacId, true
		}
	}

	func deleteVac (vacId int) (error) {
			client,ctx,err:=gprcConnect()
			if err != nil {
				fmt.Println("RPC connect error")
				return err
			}
			result, err := client.Delete(ctx, &proto.DeleteRequest{DeleteTarget: int32(vacId)})
			if err != nil {
				fmt.Println("RPC execution error")
				return err
			}
			fmt.Println(result.DeleteResult)
			return err
		}
	



func searchDialog() (string, bool) {
	fmt.Println("Введите название вакансии частично или полностью, либо наберите \"назад\" для выхода в предыдущее меню")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			if scanner.Text() == "назад" {
				return "", false
			}
			searchKey := scanner.Text()
			return searchKey, true
		}
	}
}

func loadVacs(qry string) ([]*proto.VacancyStruct, error) {
client,ctx,err:=gprcConnect()
if err != nil {
	fmt.Println("RPC error")
	return nil, err
}
	resp, err := client.Read(ctx, &proto.ReadRequest{ReadQuery: qry})
	if err != nil {
		fmt.Println("RPC error")
		return nil, err
	}
	return resp.ReadResult, err
}

func showVacs(resSlice []*proto.VacancyStruct) error {
	fmt.Println("showing")
	var counter int
	var err error
	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for i, line := range resSlice {
		counter++
		if i == 0 {

			_, err = fmt.Fprintln(w, "\tID\tНазвание вакансии\tКлючевые навыки\tОписание вакансии\tЗарплата\tТип работы\t")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(w, "\t--\t-----------------\t------------------------------------------\t-----------------------------------------------------------------\t--------\t----------\t")
			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintf(w, "\t%v\t%v\t%v\t%v\t%v\t%v\t\n", line.ID, line.VacName, line.KeySkills, line.VacDesc, line.Salary, line.JobType)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, "\t--\t-----------------\t------------------------------------------\t-----------------------------------------------------------------\t--------\t----------\t")
		if err != nil {
			return err
		}
	}
	w.Flush()
	if counter == 0 {
		fmt.Println("\nПохоже, по такому запросу в базе ничего не нашлось. Попробуйте изменить запрос")
		fmt.Println("----------------------------------------")
		return err
	}

	return err
}
func insertDialog() (proto.VacancyStruct, bool) {
	var result proto.VacancyStruct
	fmt.Println("введите соответствующие значения строк, разделяя их знаком \"/\": ")
	fmt.Println("название вакансии, ключевые навыки, описание вакансии, зарплата, и код типа работы: 1 для работы в офисе, 2 для удаленной работы и 3 для гибридной формы работы")
	fmt.Println("Например: \"Охранник/Решать конфликтные ситуации, обращаться с оружием, разгадывать сканворды/Человек, который следит за порядком в офисном здании/50000/1\"")
	fmt.Println("или наберите \"назад\" для выхода в предыдущее меню")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			if scanner.Text() == "назад" {
				return proto.VacancyStruct{}, false
			}

			queryString := strings.Split(scanner.Text(), "/")
			if len(queryString) != 5 {
				fmt.Println("Неверное количество аргументов, повторите ввод")
				continue
			}
			result.VacName = queryString[0]
			result.KeySkills = queryString[1]
			result.VacDesc = queryString[2]
			intSal, err := strconv.Atoi(queryString[3])

			if err != nil {
				fmt.Println("Ошибка ввода данных в поле \"Зарплата\", повторите ввод")
				continue
			}
			result.Salary = int32(intSal)
			intCode, err := strconv.Atoi(queryString[4])
			if err != nil {
				fmt.Println("Ошибка ввода данных в поле \"код типа работы\", повторите ввод")
				continue
			}
			result.JobCode = int32(intCode)
			if result.JobCode > 3 || result.JobCode < 1 {
				fmt.Println("Код работы может быть только следующих значений: 1 для работы в офисе, 2 для удаленной работы и 3 для гибридной формы работы. Ввод других значений не допускается")
				continue
			}
		}

		return result, true

	}
}
func insert(q proto.VacancyStruct) error {
	client,ctx,err:=gprcConnect()
	if err != nil {
		fmt.Println("RPC error")
		return err
	}
	result, err := client.Insert(ctx, &proto.InsertRequest{NewVac: &proto.VacancyStruct{VacName: q.VacName, KeySkills: q.KeySkills, VacDesc: q.VacDesc, Salary: int32(q.Salary), JobCode: int32(q.JobCode)}})
	if err != nil {

		return err
	}
	fmt.Println(result)
	return nil
}

func mainDialog() error {
	fmt.Println(mainMenu)
	scanner := bufio.NewScanner(os.Stdin)
OuterLoop:
	for {
		// Print a prompt
		fmt.Print("> ")

		// Scan for input
		if scanner.Scan() {
			switch {
			case scanner.Text() == "посмотреть":
				res, err := loadVacs("")
				if err != nil {
					fmt.Println("Ошибка обращения к серверу", err)
					return err
				}
				err = showVacs(res)
				if err != nil {
					fmt.Println("Ошибка обращения к серверу", err)
					return err
				}
				fmt.Println(mainMenu)
			case scanner.Text() == "найти":
				keyWord, proceed := searchDialog()
				if proceed {
					res, err := loadVacs(keyWord)
					if err != nil {
						fmt.Println("Ошибка обращения к серверу", err)
						return err
					}
					err = showVacs(res)
					if err != nil {
						fmt.Println("Ошибка обращения к серверу", err)
						return err
					}
				}
				fmt.Println(mainMenu)
			case scanner.Text() == "добавить":
				query, proceed := insertDialog()
				if proceed {
					err := insert(query)
					if err != nil {
						fmt.Println("Ошибка внесения данных в таблицу, попробуйте еще раз", err)
						return err
					}
				}
				fmt.Println(mainMenu)
			case scanner.Text() == "удалить":
				query, proceed := deleteDialog()
				if proceed {
					err := deleteVac(query)
					if err != nil {
						fmt.Println("Ошибка удаления данных в таблицу, попробуйте еще раз", err)
						return err
					}
				}
				fmt.Println(mainMenu)
			case scanner.Text() == "выход":
				fmt.Println("Всего хорошего!")
				break OuterLoop
			default:
				fmt.Println("Неверно введена команда, попробуйте еще раз")
			}
		}
	}
	return nil
}
