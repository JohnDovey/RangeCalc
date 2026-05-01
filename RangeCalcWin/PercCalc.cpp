#include <iostream>
#include <string>
#include <iomanip>
#include <limits>

int main()
{
    const std::string ProgVer = "1.0.2";

    // Header
    std::cout << "+------------------------------[" << std::setw(6) << ProgVer << "]----------------------------+\n";
    std::cout << "|          Percentage Calculator - by John Dovey                     |\n";
    std::cout << "+--------------------------------------------------------------------+\n\n";

    std::cout << "Enter values as prompted (type 0 to exit)\n\n";

    try
    {
        double firstNum = GetNumber("Enter First Number");
        double secondNum = GetNumber("Enter Second Number");

        double percentage = CalculatePercentage(firstNum, secondNum);

        std::cout << "\n";
        std::cout << std::fixed << std::setprecision(2);
        std::cout << "The Percentage is " << percentage << "%\n";
    }
    catch (const std::exception& ex)
    {
        std::cout << "\nError: " << ex.what() << std::endl;
    }

    std::cout << "\nHit any key to exit...";
    std::cin.get(); // Wait for user input
    return 0;
}

// Get a valid number from user with validation loop
double GetNumber(const std::string& prompt)
{
    while (true)
    {
        std::cout << prompt << ": ";
        std::string input;
        std::getline(std::cin, input);

        // Remove whitespace
        input.erase(0, input.find_first_not_of(" \t"));
        input.erase(input.find_last_not_of(" \t") + 1);

        if (input.empty())
            continue;

        try
        {
            double value = std::stod(input);

            if (value == 0.0)
                throw std::runtime_error("Operation cancelled by user.");

            if (value < 0.0)
            {
                std::cout << "Please enter a positive number.\n";
                continue;
            }

            std::cout << "  (" << value << ")\n";
            return value;
        }
        catch (...)
        {
            std::cout << "Invalid input. Please enter a valid number.\n";
        }
    }
}

// Calculate percentage safely
double CalculatePercentage(double first, double second)
{
    if (second == 0.0)
        throw std::runtime_error("Cannot divide by zero!");

    return (first * 100.0) / second;
}
