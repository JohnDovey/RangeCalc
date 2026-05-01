#include <iostream>
#include <iomanip>
#include <string>
#include <cmath>
#include <limits>

const std::string ProgVer = "1.0.15";

void PrintHeader()
{
    std::cout << "+------------------------------[" << std::setw(6) << ProgVer << "]----------------------------+\n";
    std::cout << "|        Range Calculator - by John Dovey <dovey.john@gmail.com>     |\n";
    std::cout << "| View the code on GitHub [https://github.com/JohnDovey/RangeCalc]   |\n";
    std::cout << "+--------------------------------------------------------------------+\n\n";
}

int GetBearing(const std::string& prompt)
{
    while (true)
    {
        std::cout << prompt << ": ";
        std::string input;
        std::getline(std::cin, input);

        try
        {
            int value = std::stoi(input);

            if (value == 0)
                throw std::runtime_error("Operation cancelled.");

            if (value < 1 || value > 360)
            {
                std::cout << "Error: Bearing must be between 1 and 360 degrees.\n\n";
                continue;
            }

            std::cout << "\t(" << value << "°)\n";
            return value;
        }
        catch (...)
        {
            std::cout << "Invalid input. Please enter a number.\n\n";
        }
    }
}

int GetPositiveDistance(const std::string& prompt)
{
    while (true)
    {
        std::cout << prompt << ": ";
        std::string input;
        std::getline(std::cin, input);

        try
        {
            int value = std::stoi(input);

            if (value == 0)
                throw std::runtime_error("Operation cancelled.");

            if (value < 1)
            {
                std::cout << "Error: Distance must be greater than zero.\n\n";
                continue;
            }

            std::cout << "\t(" << value << " meters between reference points)\n";
            return value;
        }
        catch (...)
        {
            std::cout << "Invalid input. Please enter a positive number.\n\n";
        }
    }
}

double CalculateRange(int b1, int b2, int distance)
{
    int diff = std::abs(b1 - b2);

    if (diff == 0)
        throw std::runtime_error("Bearings are identical - cannot calculate range.");

    if (diff >= 90)
        throw std::runtime_error("Angle difference too large (target likely behind baseline).");

    // Corrected: tan(90° - θ) = cot(θ)
    double radians = (90.0 - diff) * M_PI / 180.0;
    return distance * std::tan(radians);
}

int main()
{
    PrintHeader();
    std::cout << "Enter the values as prompted (0 to exit)\n\n";

    try
    {
        int bearing1 = GetBearing("Bearing from Ref One to Target");
        int bearing2 = GetBearing("Bearing from Ref Two to Target");
        int refDistance = GetPositiveDistance("Distance between Ref One and Ref Two");

        double range = CalculateRange(bearing1, bearing2, refDistance);

        std::cout << "\n";
        std::cout << std::fixed << std::setprecision(2);
        std::cout << "Range to Target: " << range << " meters\n\n";
    }
    catch (const std::exception& e)
    {
        std::cout << "\nError: " << e.what() << std::endl;
    }

    std::cout << "Press any key to exit...";
    std::cin.get();
    return 0;
}
