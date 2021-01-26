#include <stdio.h>
#include <string.h>
#include <cmath>

int main()
{
	int bearing1;		  // bearing from left ref point to target
	int bearing2;		  // bearing from right ref point to target
	int refdistance;	  // distance between two ref points _ use %u for unsigned number
	double rangetotarget; // Calculated range
	int tmpAngle;

	printf("\n\n Enter the values as prompted\r\n");

	// Bearing One
	printf("Bearing from Ref One to Target: ");
	scanf_s("%i", &bearing1);
	printf("\t(%u degrees bearing Ref 1 to target) \r\n", bearing1);

	if (bearing1 > 360)
	{
		printf("Error! No more than 360 Degrees allowed\r\n");
		return 1;
	}

	if (bearing1 < 1)
	{
		printf("Error! Bearing must be greater than zero\r\n");
		return 10;
	}
	// Bearing Two
	printf("Bearing from Ref Two to Target: ");
	scanf_s("%i", &bearing2);
	printf("\t(%u degrees bearing Ref 2 to target) \r\n", bearing2);

	if (bearing2 > 360)
	{
		printf("Error! No more than 360 Degrees allowed\r\n");
		return 2;
	}
	if (bearing2 < 1)
	{
		printf("Error! Bearing must be greater than zero\r\n");
		return 20;
	}

	// Seperation Distance
	printf("Distance between Ref One and Ref Two: ");
	scanf_s("%i", &refdistance);
	printf("\t(%u meters between reference points) \r\n", refdistance);

	// d = (Tan (90 - (A -B))) x Ref
	tmpAngle = bearing1 - bearing2;
	if (tmpAngle < 0)
	{
		tmpAngle = tmpAngle * (0 - 1);
	}
	rangetotarget = (tan(90 - tmpAngle)) * refdistance;
	printf("\r\nRange to Target: %f\r\n\r\n", rangetotarget);
	return 0;
}
