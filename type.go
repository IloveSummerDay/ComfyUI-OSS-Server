package main

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type SuccessResponse struct {
	Data []OSSFile `json:"data"`
}

type SaveOSSRequest struct {
	File_name_list []string `json:"file_name_list" binding:"required"`
	Client_id      string   `json:"client_id" binding:"required"`
	Ai_server_host string   `json:"ai_server_host" binding:"required"`
	Ai_server_port string   `json:"ai_server_port" binding:"required"`
}

type OSSFile struct {
	Filename string `json:"filename"`
	OSS      string `json:"oss"`
}
