#include<stdio.h>
#include<wiringPi.h>
#define PIN 11
int main()
{
	int i,j;
	unsigned long databuf = 0;  //存放温湿度内存
	unsigned char crc = 0;
//初始化wiringPi库
	if(wiringPiSetup() == -1)
	{
		printf("error wiringPi init\n");
		return -1;
	}
    //初始化引脚的工作模式
	pinMode(PIN,OUTPUT);
//初始化总线，拉高电平（只对这个温湿度而言，其它根据引脚的功能而定）
	digitalWrite(PIN,HIGH);//init 
 
 
	while(1)
	{
		delay(3000);
		pinMode(PIN,OUTPUT); 
		digitalWrite(PIN,LOW);//拉低电平，延时25ms（t1时刻），触发DH11起始信号（叫它准备工作）
		delay(25);
		digitalWrite(PIN,HIGH);//拉高电平，延时27us（t2时刻），主机开始发起结束信号
		//delayMicroseconds(27);
		pinMode(PIN,INPUT);    //设置为输入模式
		pullUpDnControl(PIN, PUD_UP);//启用上拉电阻，将不稳定的信号设为高电平，以免干扰
		
		delayMicroseconds(27);
		if (digitalRead(PIN) == 0)   //dh11拉低电平，延时80us（t3时刻），（功能：向主机发起响应信号）
		{
			while(!digitalRead(PIN))   //等待高电平（t4时刻）
				;
                //dh11初始化成功（获得数据）
 
                //开始传输前32位温湿度数据
			for(i=0;i<32;i++)          
			{
				while(digitalRead(PIN))  //等待低电平（t5时刻）
				;
				while(!digitalRead(PIN))//等待高电平（t6时刻）高电平的响应时长来判断是1或0
				;
				delayMicroseconds(32);
				databuf *=2;            //左移一位（左高右低）
				if(digitalRead(PIN)==HIGH)//延时32us后还是高电平 说明是1
				{
					databuf++;     //是1databuf加一
				}
 
			}
            //开始传输最后八位校验位
			for(i=0;i<8;i++)
			{
				while(digitalRead(PIN))//等待低电平（类t5时刻）
					;
				while(!digitalRead(PIN))//等待高电平（类t6时刻）高电平的响应时长来判断是1或0
					;
				delayMicroseconds(32);
				crc *=2;
				if(digitalRead(PIN)==HIGH)//延时32us后还是高电平 说明是1
				{
					crc++;
				}
 
			}
		}
//发送是先发低位依次到高位，高位还是高位，低位是否是低位，具体看大端还是小端模式
        //dh11前八位是湿度整数（最低位），之后八位是小数，在之后八位是温度整数，在之后八位是温度小数，最后八位是校验位（最高位），跟我们平时的读写习惯相反，左高又低
        //字节序较高位储存在字节序较低位，所以湿度整数存储在最低位，进行与运算有1的都为1，为0的都为零
		printf("RH:%d.%d\n",(databuf>>24) & 0xff,(databuf>>16) & 0xff);
		printf("TMP:%d.%d\n",(databuf>>8) & 0xff,databuf & 0xff);
		databuf = 0;
	}
 
	return 0;
}
