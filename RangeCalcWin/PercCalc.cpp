#include <iostream>

using namespace std;

int main()
{

	char FirstNum[20]; //First Number
	char SecondNum[20]; //Second Number
	int MyAns = 0;//My Answer
	int iFirstNum = 0;
	int iSecondNum = 0;


	cout << "Enter First Number and hit ENTER " <<endl;
	cin >> FirstNum;
	iFirstNum = stoi(FirstNum);
	cout << "Enter Second Number and hit ENTER "<<endl;
	cin >> SecondNum;
	iSecondNum = stoi(SecondNum);

	MyAns = (iFirstNum * 100) / iSecondNum;

	cout << "The Percentage is "<< MyAns <<"%" <<endl;

    return 0;
}
