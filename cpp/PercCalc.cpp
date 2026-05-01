#include <iostream>
#include <string>

int main()
{
    double first, second;

    std::cout << "Enter First Number: ";
    std::cin >> first;

    std::cout << "Enter Second Number: ";
    std::cin >> second;

    if (second == 0)
    {
        std::cout << "Error: Cannot divide by zero!\n";
        return 1;
    }

    double percentage = (first * 100.0) / second;
    std::cout << "The Percentage is " << percentage << "%\n";

    return 0;
}
