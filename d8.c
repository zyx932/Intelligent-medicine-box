#include <wiringPi.h>

int main(void) {
    // 初始化wiringPi库，使用BCM GPIO编号方案
    if (wiringPiSetup() == -1) {
        // 如果初始化失败，打印错误信息并退出
        printf("设置wiringPi失败!\n");
        return 1;
    }

    // 设置GPIO 8为输出模式
    pinMode(8, OUTPUT);

    // 将GPIO 8设为低电平
    digitalWrite(8, LOW);

    // 程序执行成功
    printf("GPIO 8已设为低电平\n");

    return 0;
}
