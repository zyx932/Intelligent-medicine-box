#include <stdio.h>
#include <wiringPi.h>

#define PIN 8 // 定义要操作的GPIO引脚编号

int main() {
//    if (wiringPiSetup() == -1) {
//        printf("无法初始化wiringPi库: %s\n", strerror(errno));
//        return 1;
//    }


    wiringPiSetup();
    pinMode(PIN, OUTPUT); // 设置GPIO引脚为输出模式
    digitalWrite(PIN, HIGH); // 将GPIO引脚设为高电平

    printf("GPIO %d 已设为高电平\n", PIN);

    return 0;
}
