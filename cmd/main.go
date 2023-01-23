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
var mainMenu string = "\n*****************************\n - Если хотите посмотреть всю таблицу вакансий, наберите \"посмотреть\", \n - Если хотите найти вакансию по названию наберите \"найти\"\n - Если хотите добавить строку - наберите \"добавить\", \n - Если хотите выйти из программы, наберите \"выход\"\n*****************************\n"
var searchQry string = " SELECT vacancies.id, vacancy_name, key_skills, salary, vacancy_desc, job_types.job_type FROM vacancies JOIN job_types ON vacancies.job_type = job_types.id WHERE vacancy_name ILIKE '%"


func main() {

	err:=mainDialog()
	if err!=nil{
		fmt.Println(err)
	}



}



type vacQuery struct {
	ID int
	vacName   string
	keySkills string
	vacDesc   string
	salary    int
	jobCode   int
	jobType string
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

func loadVacs(qry string) [] vacQuery{
		ctx := context.TODO()
	conn, err := grpc.Dial("172.27.148.130:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("No calculus")
	}
	defer conn.Close()
	client := proto.NewDBServerClient(conn)
	resp, err :=client.Get(ctx, &proto.GetRequest{Query: qry})
result:= strings.Split(resp.String(), "\n")
	
}

func showVacs(resSlice []vacQuery) error {
var counter int
	var err error
	const padding = 1


	



	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for i,line:= range resSlice {
counter++
			if i == 0 {
			
			_,err= fmt.Fprintln(w, "\tID\tНазвание вакансии\tКлючевые навыки\tОписание вакансии\tЗарплата\tТип работы\t")
			if err != nil {
				return err
			}
			_,err= fmt.Fprintln(w, "\t--\t-----------------\t------------------------------------------\t-----------------------------------------------------------------\t--------\t----------\t")
			if err != nil {
				return err
			}
		}
		

		_,err= fmt.Fprintf(w, "\t%v\t%v\t%v\t%v\t%v\t%v\t\n",line.ID, line.vacName, line.keySkills, line.vacDesc, line.salary, line.jobType)
		if err != nil {
			return err
		}
		_,err= fmt.Fprintln(w, "\t--\t-----------------\t------------------------------------------\t-----------------------------------------------------------------\t--------\t----------\t")
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
func insertDialog() (vacQuery, bool) {
	var result vacQuery
	var err error
	fmt.Println("введите соответствующие значения строк, разделяя их знаком \"/\": ")
	fmt.Println("название вакансии, ключевые навыки, описание вакансии, зарплата, и код типа работы: 1 для работы в офисе, 2 для удаленной работы и 3 для гибридной формы работы")
	fmt.Println("Например: \"Охранник/Решать конфликтные ситуации, обращаться с оружием, разгадывать сканворды/Человек, который следит за порядком в офисном здании/50000/1\"")
	fmt.Println("или наберите \"назад\" для выхода в предыдущее меню")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			if scanner.Text() == "назад" {
				return vacQuery{}, false
			}

			queryString := strings.Split(scanner.Text(), "/")
			if len(queryString) != 5 {
				fmt.Println("Неверное количество аргументов, повторите ввод")
				continue
			}
			result.vacName = queryString[0]
			result.keySkills = queryString[1]
			result.vacDesc = queryString[2]
			result.salary, err = strconv.Atoi(queryString[3])
			if err != nil {
				fmt.Println("Ошибка ввода данных в поле \"Зарплата\", повторите ввод")
				continue
			}
			result.jobCode, err = strconv.Atoi(queryString[4])
			if err != nil {
				fmt.Println("Ошибка ввода данных в поле \"код типа работы\", повторите ввод")
				continue
			}
			if result.jobCode > 3 || result.jobCode < 1 {
				fmt.Println("Код работы может быть только следующих значений: 1 для работы в офисе, 2 для удаленной работы и 3 для гибридной формы работы. Ввод других значений не допускается")
				continue
			}
		}

		return result, true

	}
}
func insert(q vacQuery) error {
	ctx := context.TODO()
	conn, err := grpc.Dial("172.27.148.130:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("GRPC connection error")
	}
	defer conn.Close()
	client := proto.NewDBServerClient(conn)
	result, err :=client.Put(ctx, &proto.PutRequest{VacName: q.vacName, KeySkills: q.keySkills, VacDesc:q.vacDesc, Salary: int32(q.salary), JobCode: int32(q.jobCode)})
	fmt.Println(result)

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
				res,err := loadVacs(searchQry+"%'")
				if err != nil {
					fmt.Println("Ошибка обращения к базе данных", err)
					return err
				}
				err=showVacs(res)
				if err != nil {
					fmt.Println("Ошибка обращения к базе данных", err)
					return err
				}
				fmt.Println(mainMenu)
			case scanner.Text() == "найти":
				keyWord, proceed := searchDialog()
				if proceed {
					res,err := loadVacs(searchQry+keyWord+"%'")
				if err != nil {
					fmt.Println("Ошибка обращения к базе данных", err)
					return err
				}
				err=showVacs(res)
				if err != nil {
					fmt.Println("Ошибка обращения к базе данных", err)
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
