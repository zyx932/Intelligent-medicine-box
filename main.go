package main

import (
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var cmd_play *exec.Cmd
var time_med []string
var med []string
var temp_set float64=20.0
var title string = "周煜昕的智能药箱"
var temperature float64
var humidity float64
type intoIt struct{
	TT string
	Title string
	Med string
}

func sub(hour int,minute int)(int){
	now := time.Now()
	// 设置目标时间
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	// 计算差
	d := float64(target.Sub(now))
	var dd float64=1e-09
	d=dd*d
	// 输出结果   
	return int(d)
}

func set_high() string {
	cmd := exec.Command("sudo", "./g8")
	out, err_ := cmd.Output()
	if err_ != nil {
		log.Println(err_)
	}
	log.Println("pin8 Set to high")
	return string(out)
}

func set_low() string {
	cmd := exec.Command("sudo", "./d8")
	out, err_ := cmd.Output()
	if err_ != nil {
		log.Println(err_)
	}
	log.Println("pin8 Set to low")
	return string(out)
}

func handle_set_high(w http.ResponseWriter, r *http.Request) {
	log.Println("handle Set high")
	output := set_high()
	w.Write([]byte(output))
}

func handle_set_low(w http.ResponseWriter, r *http.Request) {
	log.Println("handle Set low")
	output := set_low()
	w.Write([]byte(output))
}

func main_page(w http.ResponseWriter, r *http.Request) {
	log.Println("handle main page")
	t, _ := template.ParseFiles("./html/index.html")
	Tt:=strconv.FormatFloat(temp_set, 'f', -1, 64)
	var med_str string
	var i int
	var l int = len(time_med)
	for i<l{
		med_str=med_str+"在 "+time_med[i]+" 喝 "+med[i]+"  "
		i++
	}
	into_:=intoIt{TT:Tt,Title:title,Med:med_str}
	t.Execute(w, into_)
}

func process_data(data string) (float64, float64) {
	datas := strings.Split(data, "\n")
	log.Println("processing")
	log.Println(datas)
	var temp_str, hum_str string
	temp_str = strings.Split(datas[5], ":")[1]
	hum_str = strings.Split(datas[4], ":")[1]
	temp, _ := strconv.ParseFloat(temp_str, 64)
	hum, _ := strconv.ParseFloat(hum_str, 64)
	return temp, hum
}

func get_temperature() {
	log.Println("start getting temp")
	data := make([]byte, 1024)
	cmd := exec.Command("/home/zyx/code/bin/t3")
	ReaderAndCloser, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	cmd = exec.Command("./t3")
	ReaderAndCloser, err = cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	cmd.Start()
	data = make([]byte, 1024)
	time.Sleep(time.Second * 4)
	defer ReaderAndCloser.Close()
	defer log.Println("Getting temp ended")
	ReaderAndCloser.Read(data)
	temp, hum := process_data(string(data))
	if temp>1 && temp<40 && hum>10 && hum<88{
		temperature=temp
		humidity=hum
		if temp >= temp_set{
			set_high()
			log.Println(temp_set)
			log.Println("***Auto to set high***")
		}else{
			set_low()
			log.Println("***Auto to set low***")
		}
	}
}

func handle_get_temperature(w http.ResponseWriter, r *http.Request) {
	log.Println("handle get temperature")
	w.Write([]byte(strconv.FormatFloat(temperature, 'f', 2, 64)))
}

func handle_get_humidity(w http.ResponseWriter, r *http.Request) {
	log.Println("handle get humidity")
	w.Write([]byte(strconv.FormatFloat(humidity, 'f', 2, 64)))
}

func setting(w http.ResponseWriter,r *http.Request){
	temp,_:=template.ParseFiles("./html/setting.html")
	temp.Execute(w,title)
}

func set_tmp(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	var err error
	temp_set,err=strconv.ParseFloat(r.Form["targetTemp"][0],64)
	if err != nil{
		log.Println(err)
	}
	title=r.Form["title"][0]
}

func seting(w http.ResponseWriter, r *http.Request){
	log.Println("set time")
	r.ParseForm()
	data:=r.Form["data"][0]
	output:="set "+data+" ok"
	every:=strings.Split(data,"!")
	var i int=0
	var block string
	ii:=len(every)
	if len(every) <= 1{
		w.Write([]byte("format error"))
	}else{
		for i<ii{
			block=every[i]
			time_med=append(time_med,strings.Split(block,",")[0])
			med=append(med,strings.Split(block,",")[1])
			i++
		}
		w.Write([]byte(output))
	}
}

func help(w http.ResponseWriter,r *http.Request){
        temp,_:=template.ParseFiles("./html/help.html")
        temp.Execute(w,nil)
}

func edit(w http.ResponseWriter,r *http.Request){
	temp,_:=template.ParseFiles("./html/tj.html")
	temp.Execute(w,nil)
}



func get_() {
	for true {
		go get_temperature()
		time.Sleep(4 * time.Second)
	}
}

func tixing(){
	var i int
	l:=len(time_med)
	var block string
	var hour int
	var minute int
	for l>i{
		block=time_med[i]
		hour,_=strconv.Atoi(strings.Split(block,":")[0])
		minute,_=strconv.Atoi(strings.Split(block,":")[1])
		log.Println("###:time",sub(hour,minute))
		if sub(hour,minute) > -5 && sub(hour,minute)<5{
		//	cmd_play.Start()
			log.Println("*************play start************")
			op,_:=cmd_play.Output()
			log.Println(op)
			time.Sleep(10*time.Second)
		}
		i++
	}
}

func start_tixing(){
	for true{
		time.Sleep(time.Second)
		tixing()
	}
}

func main() {
	cmd_play=exec.Command("./play.sh")
	log.Println(sub(20,30))
	log.Println("Starting server")
	http.HandleFunc("/main",main_page)
	http.HandleFunc("/z",handle_get_humidity)
	http.HandleFunc("/x", handle_get_temperature)
	//SETTING
	http.HandleFunc("/editor",edit)
	http.HandleFunc("/settime",seting)
	http.HandleFunc("/set",setting)
	http.HandleFunc("/help",help)
	http.HandleFunc("/setitbyurl",set_tmp)
	defer log.Println("quit")
	go start_tixing()
	go get_()
	http.ListenAndServe(":8888", nil)
}
