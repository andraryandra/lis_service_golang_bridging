package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-playground/validator/v10"
    "github.com/joho/godotenv"
)

var validate *validator.Validate

func initLogger() {
    if _, err := os.Stat("logs"); os.IsNotExist(err) {
        err := os.Mkdir("logs", 0755)
        if err != nil {
            log.Fatalf("Failed to create logs directory: %v", err)
        }
    }

    logFile := filepath.Join("logs", fmt.Sprintf("log_%s.txt", time.Now().Format("2006-01-02")))
    file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }

    multiWriter := io.MultiWriter(file, os.Stdout)
    log.SetOutput(multiWriter)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// type ICD10 struct {
//     Type string `json:"type,omitempty"`
//     Code string `json:"code,omitempty"`
//     Text string `json:"text,omitempty"`
// }

type DetailRujukan struct {
    IdDokter       *string `json:"id_dokter" validate:"required"`
    NamaDokter     *string `json:"nama_dokter" validate:"required"`
    PhoneDokter    *string `json:"phone_dokter,omitempty"`
    IhsDokter      *string `json:"ihs_dokter,omitempty"`
    Nik            *string `json:"nik,omitempty"`
    IdWard         *string `json:"id_ward,omitempty"`
    Ward           *string `json:"ward,omitempty"`
    IdFasilitas    *string `json:"id_fasilitas,omitempty"`
    Fasilitas      *string `json:"fasilitas,omitempty"`
    Alamat         *string `json:"alamat,omitempty"`
    NoTelpPerujuk  *string `json:"no_telp_perujuk,omitempty"`
    IdInstalasi    *string `json:"id_instalasi,omitempty"`
}

type Order struct {
    IdTest   *string `json:"id_test" validate:"required"`
    NamaTest *string `json:"nama_test,omitempty"`
}

type SimrsRequest struct {
    NoPendaftaran     *string         `json:"no_pendaftaran" validate:"required"`
    NoRM              *string         `json:"no_rm" validate:"required"`
    NoOrder           *string         `json:"no_order" validate:"required"`
    NamaPasien        *string         `json:"nama_pasien" validate:"required"`
    TempatLahir       *string         `json:"tempat_lahir"`
    TglLahir          *string         `json:"tgl_lahir,omitempty"`
    JK                *string         `json:"jk" validate:"required"`
    Alamat            *string         `json:"alamat,omitempty"`
    Phone             *string         `json:"phone,omitempty"`
    NIK               *string         `json:"nik,omitempty"`
    IdJenisPasien     *string         `json:"id_jenis_pasien,omitempty"`
    JenisPasien       *string         `json:"jenis_pasien,omitempty"`
    IdPenjamin        *string         `json:"id_penjamin,omitempty"`
    Penjamin          *string         `json:"penjamin,omitempty"`
    RujukanAsal       *string         `json:"rujukan_asal" validate:"required"`
    DetailRujukan     []*DetailRujukan `json:"detail_rujukan" validate:"required,dive"`
    Cito              *string         `json:"cito" validate:"required"`
    Diagnose          *string         `json:"diagnose,omitempty"`
    RegistrationDate  *string         `json:"registration_date,omitempty"`
    Email             *string         `json:"email,omitempty"`
    Region            *string         `json:"region,omitempty"`
    CodeClinic        *string         `json:"code_clinic,omitempty"`
    ClinicName        *string         `json:"clinic_name,omitempty"`
    OrderNote         *string         `json:"order_note,omitempty"`
    KategoriKunjungan *string         `json:"kategori_kunjungan,omitempty"`
    ICD10             *interface{}    `json:"icd10,omitempty"`
    Order             []*Order        `json:"order" validate:"required,dive"`
}

type LISRequest struct {
    NoPendaftaran     *string         `json:"no_pendaftaran"`
    NoRM              *string         `json:"no_rm"`
    NoOrder           *string         `json:"no_order"`
    NamaPasien        *string         `json:"nama_pasien"`
    TempatLahir       *string         `json:"tempat_lahir"`
    TglLahir          *string         `json:"tgl_lahir"`
    JK                *string         `json:"jk"`
    Alamat            *string         `json:"alamat"`
    Phone             *string         `json:"phone,omitempty"`
    NIK               *string         `json:"nik"`
    IdJenisPasien     *string         `json:"id_jenis_pasien"`
    JenisPasien       *string         `json:"jenis_pasien,omitempty"`
    IdPenjamin        *string         `json:"id_penjamin"`
    Penjamin          *string         `json:"penjamin,omitempty"`
    RujukanAsal       *string         `json:"rujukan_asal"`
    DetailRujukan     []*DetailRujukan `json:"detail_rujukan"`
    Cito              *bool           `json:"cito"`
    Diagnose          *string         `json:"diagnose,omitempty"`
    RegistrationDate  *string         `json:"registration_date,omitempty"`
    Email             *string         `json:"email,omitempty"`
    Region            *string         `json:"region,omitempty"`
    CodeClinic        *string         `json:"code_clinic,omitempty"`
    ClinicName        *string         `json:"clinic_name,omitempty"`
    OrderNote         *string         `json:"order_note,omitempty"`
    KategoriKunjungan *string         `json:"kategori_kunjungan,omitempty"`
    ICD10			  *interface{}     `json:"icd10"`
    Order             []*Order        `json:"order"`
}

func transformToLIS(simrs *SimrsRequest) *LISRequest {
    var idJenisPasien, jenisPasien, idPenjamin, penjamin string

    if simrs.IdJenisPasien == nil || *simrs.IdJenisPasien == "UMUM" {
        idJenisPasien = "1"
        jenisPasien = "UMUM"
        idPenjamin = "UMUM"
        penjamin = " "
    } else {
        idJenisPasien = "2"
        jenisPasien = "ASURANSI"
        idPenjamin = *simrs.IdJenisPasien
        penjamin = *simrs.JenisPasien
    }

    // Transform Cito to bool
    var cito bool
    if simrs.Cito != nil && (*simrs.Cito == "1" || *simrs.Cito == "true") {
        cito = true
    } else {
        cito = false
    }

    var icd10 interface{}
    if simrs.ICD10 != nil {
        icd10 = simrs.ICD10
    } else {
        icd10 = []interface{}{}
    }

    // Map SimrsRequest to LISRequest
    return &LISRequest{
        NoPendaftaran:     simrs.NoPendaftaran,
        NoRM:              simrs.NoRM,
        NoOrder:           simrs.NoOrder,
        NamaPasien:        simrs.NamaPasien,
        TempatLahir:       simrs.TempatLahir,
        TglLahir:          simrs.TglLahir,
        JK:                simrs.JK,
        Alamat:            simrs.Alamat,
        Phone:             simrs.Phone,
        NIK:               simrs.NIK,
        IdJenisPasien:     &idJenisPasien,
        JenisPasien:       &jenisPasien,
        IdPenjamin:        &idPenjamin,
        Penjamin:          &penjamin,
        RujukanAsal:       simrs.RujukanAsal,
        DetailRujukan:     simrs.DetailRujukan,
        Cito:              &cito,
        Diagnose:          simrs.Diagnose,
        RegistrationDate:  simrs.RegistrationDate,
        Email:             simrs.Email,
        Region:            simrs.Region,
        CodeClinic:        simrs.CodeClinic,
        ClinicName:        simrs.ClinicName,
        OrderNote:         simrs.OrderNote,
        KategoriKunjungan: simrs.KategoriKunjungan,
        ICD10:             &icd10,
        Order:             simrs.Order,
    }
}


type SendToLisBridgingIn struct {
    request *LISRequest
    xSign   string
    xCons   string
}

func sendToLISBridging(payload SendToLisBridgingIn) (map[string]interface{}, error) {
    url := os.Getenv("LIS_BRIDGING")

    jsonData, err := json.Marshal(payload.request)
    var respBody map[string]interface{}

    if err != nil {
        log.Printf("Error marshalling request: %v", err)
        return respBody, err
    }

    log.Printf("Sending data to LIS bridging: %s", jsonData)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Error creating request: %v", err)
        return respBody, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-sign", payload.xSign)
    req.Header.Set("x-cons", payload.xCons)

    log.Printf("Request headers: %v", req.Header)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending request: %v", err)
        return respBody, err
    }
    defer resp.Body.Close()

    // Log the response from the API
    err = json.NewDecoder(resp.Body).Decode(&respBody)
    if err != nil {
        log.Printf("Error decoding response: %v", err)
        return respBody, err
    }
    log.Printf("Response from LIS bridging: %s", resp.Status)
    log.Printf("Response body: %v", respBody)

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        log.Printf("Failed to send data to LIS bridging, status code: %d", resp.StatusCode)
        return respBody, fmt.Errorf("failed to send data to LIS bridging, status code: %d", resp.StatusCode)
    }

    log.Printf("Successfully sent data to LIS bridging")
    return respBody, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    var simrsReq SimrsRequest

    // Log the incoming request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        log.Printf("Failed to read request body: %v", err)
        return
    }
    log.Printf("Incoming request body: %s", body)

    // Decode the request body into the SimrsRequest struct
    err = json.Unmarshal(body, &simrsReq)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        log.Printf("Invalid request payload: %v", err)
        return
    }

    // Validate the request payload
    err = validate.Struct(simrsReq)
    if err != nil {
        validationErrors := err.(validator.ValidationErrors)
        http.Error(w, fmt.Sprintf("Validation error: %v", validationErrors), http.StatusBadRequest)
        log.Printf("Validation error: %v", validationErrors)
        return
    }

    // Get x-sign and x-cons from headers or .env
    xSign := r.Header.Get("x-sign")
    if xSign == "" {
        xSign = os.Getenv("X_SIGN")
    }

    xCons := r.Header.Get("x-cons")
    if xCons == "" {
        xCons = os.Getenv("X_CONS")
    }

    // Log the values of x-sign and x-cons
    log.Printf("x-sign: %s", xSign)
    log.Printf("x-cons: %s", xCons)

    // Transform SIMRS DTO to LIS DTO
    lisReq := transformToLIS(&simrsReq)

    // Send the transformed data to LIS bridging
    payload := SendToLisBridgingIn{
        request: lisReq,
        xSign:   xSign,
        xCons:   xCons,
    }

    respBody, err := sendToLISBridging(payload)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error sending data to LIS bridging: %v", err), http.StatusInternalServerError)
        log.Printf("Error sending data to LIS bridging: %v", err)
        return
    }

    // Respond with the transformed data
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(respBody)
    log.Printf("Successfully processed request and sent response")
    // Log the response status
    log.Printf("Response status: %d", http.StatusOK)
    log.Printf("\n------------------------------------------------------------------------------------------------------------------\n")
}

func main() {
    // Initialize the logger
    initLogger()

    // Initialize the validator
    validate = validator.New()

    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World!"))
    })

    r.Post("/api/v1/saveOrder", handleRequest)

    httpServer := &http.Server{
        Addr:    ":8111",
        Handler: r,
    }

    log.Printf("Starting server on port %s\n", "8111")
    log.Fatal(httpServer.ListenAndServe())
}