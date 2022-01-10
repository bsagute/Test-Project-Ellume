package main

import (
	"TestProject/Server/GolangServer/collection"
	"TestProject/Server/GolangServer/confighelper"
	"TestProject/Server/GolangServer/dbhelper"
	"TestProject/Server/GolangServer/model"
	"database/sql"
	"encoding/json"
	"sync"

	_ "github.com/lib/pq"
	// "github.com/prometheus/common/log"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/tidwall/gjson"
)

// Postgres - this var will be used in the rest of code to execute DB operations
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var wg sync.WaitGroup

func main() {

	//!GO Routine DB connection Call
	// ConnectDBUsingGoRoutine()
	////!

	confighelper.InitViper()
	e := echo.New()
	//GO ROUTINE CONNECTING TO DIFFERENT DATABASE
	//Fetch the ellume branch details / City List
	e.GET("/GetEllumeBranchCityListService", GetEllumeBranchCityListService)
	//Fetch the reference material list
	e.GET("/GetMasterURLs", GetMasterURLs)
	//Fetching details using menue & City
	e.POST("/GetSolarPanalDistributerListService", GetSolarPanalDistributerListService)
	//Add the new Clinet into the system
	e.POST("/AddSolarClientDetails", AddSolarClientDetails)
	//Will Add the Employee Details into the database
	e.POST("/AddEllumeEmployeeDetails", AddEllumeEmployeeDetails)
	e.POST("/deleteRecordService", DeleteRecordService)
	e.Logger.Fatal(e.Start(":4000"))
}

func ConnectDBUsingGoRoutine() {
	wg.Add(2)
	go ConnectToSQL()
	go ConnectToMongo()
	wg.Wait()
	// time.Sleep(time.Millisecond * 100)
}

func ConnectToMongo() {
	defer wg.Done()
	db, ctx, err := dbhelper.GetMongoClient(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), confighelper.GetConfig("DBNAME"), confighelper.GetConfig("USERNAME"), confighelper.GetConfig("PASSWORD"), false)
	if err != nil {
		fmt.Println("Log DB COnnection Error")
		log.Print("Error While Connecting To MongoDB::", err)
		// return false, err
	}
	fmt.Println("GO ROUTINE CONNECTED TO MONGO ")
	fmt.Println("db, ctx,", db, ctx)
}

//COnnect to SQL
func ConnectToSQL() {
	defer wg.Done()
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("GO ROUTINE CONNECTED TO SQL ")
	fmt.Println(" SQL Connected!")
}

//Check error if DB not connecting
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

//Read the URL and retun the json string
func GetRequestBodyJson(c echo.Context) (body gjson.Result, err error) {
	bb, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("Error Reading Request Body")
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(bb), nil
}

//to delete the input login Id details PERMENTLY
func DeleteRecordService(c echo.Context) error {
	body, _ := GetRequestBodyJson(c)
	loginId := body.Get("loginId").String()

	_, serviceCallError := DeleteRecordDAO(loginId)
	if serviceCallError != nil {
		fmt.Println("Service Update Error")
		return c.JSON(http.StatusInternalServerError, serviceCallError)
	}
	return c.JSON(http.StatusOK, true)
}

//Update the user details
func AddSolarClientDetails(c echo.Context) error {

	ClientProfile := model.SolarClientModel{}
	bindError := c.Bind(&ClientProfile)

	if bindError != nil {
		fmt.Println("BIND ERROR")
		return c.JSON(http.StatusInternalServerError, bindError)
	}
	_, serviceCallError := AddSolarClientDetailsService(ClientProfile)

	if serviceCallError != nil {
		fmt.Println("Service Update Error")
		return c.JSON(http.StatusInternalServerError, serviceCallError)
	}
	return c.JSON(http.StatusOK, "Solar Client Record Inserted Successfully ")
}

//Update the user details
func AddEllumeEmployeeDetails(c echo.Context) error {

	ValunteerProfile := model.EllumeEmployeeModel{}
	bindError := c.Bind(&ValunteerProfile)

	if bindError != nil {
		fmt.Println("BIND ERROR")
		return c.JSON(http.StatusInternalServerError, bindError)
	}
	_, serviceCallError := AddEllumeEmployeeDetailsService(ValunteerProfile)

	if serviceCallError != nil {
		fmt.Println("Service Update Error")
		return c.JSON(http.StatusInternalServerError, serviceCallError)
	}
	return c.JSON(http.StatusOK, "Ellume Employee Inserted Successfully ")
}

//Update the user details
func GetSolarPanalDistributerListService(c echo.Context) error {
	body, _ := GetRequestBodyJson(c)
	selectedCity := body.Get("selectedCity").String()
	selectedMenue := body.Get("selectedMenue").String()
	fmt.Println(" sssssss", selectedCity, selectedMenue)
	cityBranchDetailsList := gjson.Result{}

	cityBranchDetailsList, serviceCallError := GetCityBranchListService(selectedCity, selectedMenue)
	if serviceCallError != nil {
		fmt.Println("Service Update Error")
		return c.JSON(http.StatusInternalServerError, serviceCallError)
	}
	return c.JSON(http.StatusOK, cityBranchDetailsList.Value())
}

func GetEllumeBranchCityListService(c echo.Context) error {

	cityMasterList := gjson.Result{}

	cityMasterList, serviceCallError := GetAllRecordListService()
	if serviceCallError != nil {
		fmt.Println("Service Update Error")
		return c.JSON(http.StatusInternalServerError, serviceCallError)
	}
	return c.JSON(http.StatusOK, cityMasterList.Value())
}

func GetMasterURLs(c echo.Context) error {

	UrlsMasterList := gjson.Result{}

	UrlsMasterList, serviceCallError := GetMasterURLService()
	if serviceCallError != nil {
		fmt.Println("Service Update Error")
		return c.JSON(http.StatusInternalServerError, serviceCallError)
	}
	return c.JSON(http.StatusOK, UrlsMasterList.Value())
}

// This method get all  record  mapping List by calling DAO method.
func GetAllRecordListService() (gjson.Result, error) {
	return GetAllRecordListDAO()
}

// This method get all  record  mapping List by calling DAO method.
func GetMasterURLService() (gjson.Result, error) {
	return GetMasterURLDAO()
}

// This method get all  record  mapping List by calling DAO method.
func GetCityBranchListService(selectedCity, selectedMenue string) (gjson.Result, error) {
	return GetEllumeBranchListDAO(selectedCity, selectedMenue)
}

//This Method Get all  recordList, Retrives data from MongoDB.
func GetAllRecordListDAO() (gjson.Result, error) {
	fmt.Println("CALLINGG")
	//TODO: Test comment
	// db, ctx, err := dbhelper.GetMongoClient("localhost", "27017", "sampleTestDatabase", "", "", false)
	fmt.Println("ssss", confighelper.GetConfig("DBNAME"))
	db, ctx, err := dbhelper.GetMongoClient(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), confighelper.GetConfig("DBNAME"), confighelper.GetConfig("USERNAME"), confighelper.GetConfig("PASSWORD"), false)
	if err != nil {
		fmt.Println("Log DB COnnection Error: GetAllRecordListDAO()")
		log.Print("Error While Connecting To MongoDB::", err)
		return gjson.Result{}, err
	}

	selector := bson.M{"isDeleted": false}
	fmt.Println("selector ", selector)
	collection := db.Database(confighelper.GetConfig("DBNAME")).Collection(collection.M_CITY)
	// ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// res, err := collection.InsertOne(ctx, bson.M{
	// 	"name": "AAAAAAA",
	// })
	res, err := collection.Find(ctx, bson.M{})

	fmt.Println("ID ", err)
	// result := res.All(ctx, res)
	// fmt.Println("result", result)
	fmt.Println("ID ", res)
	// cursor, err := db.Collection(collection.EMPLOYEE_PROFILE).Find(ctx, selector)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var records []bson.M
	if err = res.All(ctx, &records); err != nil {
		log.Fatal(err)
	}

	err = res.Decode(&records)
	if err != nil {
		log.Print(err)

	}
	bs, err := json.Marshal(records)
	if err != nil {
		log.Print(err)

		// continue
	}
	return gjson.ParseBytes(bs), nil
	// return gjson.Result{}, nil
}

//This Method Get all  recordList, Retrives data from MongoDB.
func GetMasterURLDAO() (gjson.Result, error) {
	db, ctx, err := dbhelper.GetMongoClient(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), confighelper.GetConfig("DBNAME"), confighelper.GetConfig("USERNAME"), confighelper.GetConfig("PASSWORD"), false)
	if err != nil {
		fmt.Println("Log DB COnnection Error: GetMasterURLDAO()")
		log.Print("Error While Connecting To MongoDB::", err)
		return gjson.Result{}, err
	}

	collection := db.Database(confighelper.GetConfig("DBNAME")).Collection(collection.M_URLS)
	res, err := collection.Find(ctx, bson.M{})

	fmt.Println("ID ", err)
	fmt.Println("ID ", res)
	var records []bson.M
	if err = res.All(ctx, &records); err != nil {
		log.Fatal(err)
	}

	err = res.Decode(&records)
	if err != nil {
		log.Print(err)

	}
	bs, err := json.Marshal(records)
	if err != nil {
		log.Print(err)

	}
	return gjson.ParseBytes(bs), nil
}

//This Method Get all  recordList, Retrives data from MongoDB.
func GetEllumeBranchListDAO(selectedCity, selectedMenue string) (gjson.Result, error) {
	db, ctx, err := dbhelper.GetMongoClient(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), confighelper.GetConfig("DBNAME"), confighelper.GetConfig("USERNAME"), confighelper.GetConfig("PASSWORD"), false)
	if err != nil {
		fmt.Println("Log DB COnnection Error: GetEllumeBranchListDAO()")
		log.Print("Error While Connecting To MongoDB::", err)
		return gjson.Result{}, err
	}

	selector := bson.M{"cityId": selectedCity, "menueId": selectedMenue}
	fmt.Println("selector ", selector)
	collection := db.Database(confighelper.GetConfig("DBNAME")).Collection(collection.CITY_BRANCH_DETAILS)

	res, _ := collection.Find(ctx, selector)

	var records []bson.M
	if err = res.All(ctx, &records); err != nil {
		log.Fatal(err)
	}

	err = res.Decode(&records)
	if err != nil {
		log.Print(err)

	}
	bs, err := json.Marshal(records)
	if err != nil {
		log.Print(err)

		// continue
	}
	return gjson.ParseBytes(bs), nil
	// return gjson.Result{}, nil
}

// This method delete personal information by calling DAO method.
func DeleteRecordDAO(loginId string) (bool, error) {

	flag, updateServiceError := DeleteDAO(loginId)
	if updateServiceError != nil {
		fmt.Println(" UpdateService Error")
		return false, updateServiceError
	}
	return flag, nil
}

// This method update personal information by calling DAO method.
func AddSolarClientDetailsService(ClientProfile model.SolarClientModel) (bool, error) {

	flag, AddSolarClientDetailsServiceError := AddSolarClientDetailsDAO(ClientProfile)
	if AddSolarClientDetailsServiceError != nil {
		return false, AddSolarClientDetailsServiceError
	}
	return flag, nil
}

// This method update personal information by calling DAO method.
func AddEllumeEmployeeDetailsService(ClientProfile model.EllumeEmployeeModel) (bool, error) {

	flag, AddValunteerDetailsServiceError := AddEllumeEmployeeDetailsDAO(ClientProfile)
	if AddValunteerDetailsServiceError != nil {
		return false, AddValunteerDetailsServiceError
	}
	return flag, nil
}

// This Method Delete personal information, Update data in MongoDB.
func DeleteDAO(loginId string) (bool, error) {
	db, ctx, err := dbhelper.GetMongoClient(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), confighelper.GetConfig("DBNAME"), confighelper.GetConfig("USERNAME"), confighelper.GetConfig("PASSWORD"), false)

	// db, ctx, err := dbhelper.GetMongoClient(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), "sampleTestDatabase", "", "", false)
	if err != nil {
		fmt.Println("Log DB COnnection Error")
		log.Print("Error While Connecting To MongoDB::", err)
		return false, err
	}
	selector := bson.M{"loginId": loginId}
	fmt.Println("", selector, db, ctx)
	// _, err = db.Collection(collection.EMPLOYEE_PROFILE).DeleteOne(ctx, selector)
	// if err != nil {
	// 	log.Fatal("Error While Deleting Record", err)
	// }
	return true, nil
}

// This Method update personal information, Update data in MongoDB.
func AddSolarClientDetailsDAO(templateModelObj model.SolarClientModel) (bool, error) {
	db, ctx, err := dbhelper.GetMongoClient(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), confighelper.GetConfig("DBNAME"), confighelper.GetConfig("USERNAME"), confighelper.GetConfig("PASSWORD"), false)
	if err != nil {
		fmt.Println("Log DB COnnection Error")
		log.Print("Error While Connecting To MongoDB::", err)
		return false, err
	}
	collection := db.Database(confighelper.GetConfig("DBNAME")).Collection(collection.SOLAR_CLIENT_DETAILS)
	res, err := collection.InsertOne(ctx, bson.M{
		"FirstName": templateModelObj.FirstName,
		"Age":       templateModelObj.Age,

		"MiddleName": templateModelObj.MiddleName,
		"LastName":   templateModelObj.LastName})
	fmt.Println("Inserted ID ", res.InsertedID)

	return true, err
}

// This Method update personal information, Update data in MongoDB.
func AddEllumeEmployeeDetailsDAO(templateModelObj model.EllumeEmployeeModel) (bool, error) {
	db, ctx, err := dbhelper.GetMongoClient(confighelper.GetConfig("DBIP"), confighelper.GetConfig("PORT"), confighelper.GetConfig("DBNAME"), confighelper.GetConfig("USERNAME"), confighelper.GetConfig("PASSWORD"), false)
	if err != nil {
		fmt.Println("Log DB COnnection Error")
		log.Print("Error While Connecting To MongoDB::", err)
		return false, err
	}
	collection := db.Database(confighelper.GetConfig("DBNAME")).Collection(collection.ELLUME_EMPLOYEE_DETAILS)
	res, err := collection.InsertOne(ctx, bson.M{
		"FirstName":  templateModelObj.FirstName,
		"Age":        templateModelObj.Age,
		"Remark":     templateModelObj.Remark,
		"BloodGroup": templateModelObj.BloodGroup,
		"MobileNo":   templateModelObj.MobileNo,
		"EmailId":    templateModelObj.EmailId,
		"FullName":   templateModelObj.FullName,
		"LastName":   templateModelObj.LastName,
	})
	fmt.Println("Inserted ID ", res.InsertedID)

	return true, err
}
