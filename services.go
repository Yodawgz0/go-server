package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
)

type UserData struct {
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved GET reuqest from %s for %s", r.RemoteAddr, r.URL.Path)

	responseMessage := fmt.Sprintf("Hello, this is your Go server!")
	fmt.Fprint(w, responseMessage)
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "ashley",
		Password: "bazzi",
	}
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS vehicle_census WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}").Exec()
	if err != nil {
		log.Fatal(err)
	}
	session.Close()
	cluster.Keyspace = "vehicle_census"
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	err = session.Query(`CREATE TABLE IF NOT EXISTS vehicle_census (
		ID int,
		TABWEIGHT float,
		REGSTATE text,
		ACQUIREYEAR text,
		ACQUISITION int,
		AVGWEIGHT text,
		BRAKES text,
		BTYPE text,
		BUSRELATED int,
		CAB text,
		CABDAY text,
		CABHEIGHT text,
		CI_AUTOEBRAKE text,
		CI_AUTOESTEER text,
		CI_RAUTOEBRAKE text,
		CUBICINCHDISP text,
		CW_BLINDSPOT text,
		CW_FWDCOLL text,
		CW_LANEDEPART text,
		CW_PARKOBST text,
		CW_RCROSSTRAF text,
		CYLINDERS text,
		DC_ACTDRIVASST text,
		DC_ADAPCRUISE text,
		DC_LANEASST text,
		DC_PLATOON text,
		DC_VTVCOMM text,
		DEADHEADPCT text,
		DRIVEAXLES int,
		ENGREBUILD text,
		ER_COMPOWN text,
		ER_COST text,
		ER_DEALER text,
		ER_GENERAL text,
		ER_LEASING text,
		ER_OTHER text,
		ER_SELF text,
		ER_UNKLOC text,
		FE_AEROBUMP text,
		FE_AEROHOOD text,
		FE_AEROMIRROR text,
		FE_AUTOENGOFF text,
		FE_AUTOTIREINF text,
		FE_FAIRINGS text,
		FE_FTCOVER text,
		FE_GAPREDUCER text,
		FE_HYBDRIVENP text,
		FE_HYBDRIVEPL text,
		FE_IDLEREDUCE text,
		FE_LRRTIRES text,
		FE_NOSECONE text,
		FE_SIDESKIRT text,
		FTYPE text,
		FUELTYPE text,
		FUELWHERE text,
		GM_COMPOWN text,
		GM_COST text,
		GM_DEALER text,
		GM_GENERAL text,
		GM_LEASING text,
		GM_OTHER text,
		GM_SELF text,
		GVWR_CLASS text,
		HAZCARRY text,
		HAZPCT text,
		HBSTATE text,
		HBTYPE text,
		JU_CANADA text,
		JU_HOMEBASE text,
		JU_MEXICO text,
		JU_OTHERST text,
		KINDOFBUS text,
		LE_HEIGHT text,
		LE_WEIGHTBR text,
		LE_WEIGHTHI text,
		LEASECHAR text,
		LEASELENGTH text,
		LEASESTAT text,
		LF_BELOW text,
		LF_BEYOND text,
		LF_FORWARD text,
		LOADEDPCT text,
		LP_FINANCEONLY text,
		LP_FUELCONT text,
		LP_FULLMAINT text,
		LP_LICENSEPERM text,
		LP_PAYTAX text,
		LP_RECORDKEEP text,
		LTRUCKLOADPCT text,
		MILESANNL int,
		MILESLIFE int,
		MODELYEAR int,
		MONTHOPERATE int,
		MPG float,
		NUMBRAKING int,
		NUMGEARS int,
		NUMLIFTABLE int,
		OD_AHIGHBEAM text,
		OD_BACKUPCAM text,
		OD_DRIVERMON text,
		OD_HUD text,
		OD_NIGHTVIS text,
		OD_SVCAM text,
		OF_AERIAL text,
		OF_AIRCOMPRESS text,
		OF_AIRSPRING text,
		OF_AUXGEN text,
		OF_CRANE text,
		OF_DIGGERDER text,
		OF_EMERLIGHT text,
		OF_ENGINERET text,
		OF_HITCH text,
		OF_HOIST text,
		OF_INVERTER text,
		OF_LADDERRACK text,
		OF_LIFTGATE text,
		OF_MOUNTBAR text,
		OF_PARTITION text,
		OF_POWTAKEOFF text,
		OF_RAILWAYAXLE text,
		OF_REFRIGERATOR text,
		OF_SPREADER text,
		OF_TELEMATICS text,
		OF_TOOLBOX text,
		OF_WELDER text,
		OF_WINCH text,
		OF_WSBYPASS text,
		OWGTPMTANN float,
		OWGTPMTSNG float,
		PA_APARKASST text,
		PA_REMOTEPARK text,
		PAXLECONFIG text,
		PRIMCOMMACT text,
		PRIMPROD text,
		REARAXLETIRES text,
		REPLACE text,
		REPOSITIONPCT text,
		RGROUP text,
		RO_0_50 text,
		RO_101_200 text,
		RO_201_500 text,
		RO_51_100 text,
		RO_GT500 text,
		ST_ABS text,
		ST_AIRBAG text,
		ST_CRUISE text,
		ST_DRIVERCAM text,
		ST_GPS text,
		ST_GPSNAV text,
		ST_INTERNET text,
		ST_ROLLOVER text,
		TCONFIG text,
		TE_AEROREF text,
		TE_ALUMWHEEL text,
		TE_FRONTFAIRING text,
		TE_LWLANDGEAR text,
		TE_OTHER text,
		TE_REARFAIRING text,
		TE_SIDESKIRTS text,
		TE_UCAERODEV text,
		TOTLENGTH int,
		TOWCAPACITY int,
		TRANSMISSION text,
		TRIPOFFROAD int,
		TRUCKLOADPCT text,
		TTYPE text,
		VEHTYPE text,
		WEIGHOUTPCT text,
		PRIMARY KEY (ID)
	);`).Exec()

}
func handleReadGetRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved GET reuqest from %s for %s", r.RemoteAddr, r.URL.Path)

	// Create a cluster configuration and a session using the keyspace
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Keyspace = "vehicle_census"
	// other options
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Query the system_schema.tables table to list the tables in the keyspace
	iter := session.Query("SELECT vehicle_census WHERE keyspace_name = ?", "vehicle_census").Iter()
	var tableName string
	for iter.Scan(&tableName) {
		// Print the table name to the console
		fmt.Println(tableName)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	responseMessage := fmt.Sprintf("Reading of the table name is done!")
	fmt.Fprint(w, responseMessage)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received POST request from %s for %s", r.RemoteAddr, r.URL.Path)
	var userData UserData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Data received: Username - %s, Age - %d", userData.Username, userData.Age)

	// Send a response to the client
	responseMessage := fmt.Sprintf("Data received is Username: %s and Age: %d",
		userData.Username, userData.Age)
	fmt.Fprintf(w, responseMessage)
}
