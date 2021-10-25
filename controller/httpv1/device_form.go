package httpv1

type CreateDeviceForm struct{
	Name string `form:"name" json:"name" binding:"required,max=10"`
	Type string `form:"type" json:"type" binding:"required"`
	ServerIp string `form:"server_ip" json:"server_ip" binding:"required,max=15"`
	ServerPort string `form:"server_port" json:"server_port" binding:"required,max=5"`
}

