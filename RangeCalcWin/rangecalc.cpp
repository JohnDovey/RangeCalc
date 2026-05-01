#include <iostream>
#include <iomanip>
#include <string>
#include <cmath>
#include <limits>

int main()
{
    const std::string ProgVer = "1.0.14";

    // Header
    std::cout << "+------------------------------[" << std::setw(6) << ProgVer << "]----------------------------+\n";
    std::cout << "|        Range Calculator - by John Dovey <dovey.john@gmail.com>     |\n";
    std::cout << "| View the code on GitHub [https://github.com/JohnDovey/RangeCalc]   |\n";
    std::cout << "+--------------------------------------------------------------------+\n\n";

    std::cout << "Enter the values as prompted (type 0 to exit)\n\n";

    try
    {
        int bearing1 = GetBearing("Bearing from Ref One to Target");
        int bearing2 = GetBearing("Bearing from Ref Two to Target");
        int refDistance = GetPositiveInt("Distance between Ref One and Ref Two (meters)");

        double range = CalculateRange(bearing1, bearing2, refDistance);

        std::cout << "\n";
        std::cout << std::fixed << std::setprecision(2);
        std::cout << "Range to Target: " << range << " meters\n";
    }
    catch (const std::exception& ex)
    {
        std::cout << "\nError: " << ex.what() << std::endl;
    }

    std::cout << "\nHit any key to exit...";
    std::cin.get();
    return 0;
}

// Get valid bearing (1-360)
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
                throw std::runtime_error("Operation cancelled by user.");

            if (value < 1 || value > 360)
            {
                std::cout << "Error: Bearing must be between 1 and 360 degrees.\n";
                continue;
            }

            std::cout << "\t(" << value << " degrees)\n";
            return value;
        }
        catch (...)
        {
            std::cout << "Please enter a valid number.\n";
        }
    }
}

// Get valid positive distance
int GetPositiveInt(const std::string& prompt)
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
                throw std::runtime_error("Operation cancelled by user.");

            if (value < 1)
            {
                std::cout << "Error: Distance must be greater than zero.\n";
                continue;
            }

            std::cout << "\t(" << value << " meters between reference points)\n";
            return value;
        }
        catch (...)
        {
            std::cout << "Please enter a valid positive number.\n";
        }
    }
}

// Calculate range - FIXED: proper radians conversion
double CalculateRange(int bearing1, int bearing2, int refDistance)
{
    int angleDiff = std::abs(bearing1 - bearing2);

    if (angleDiff == 0)
        throw std::runtime_error("Bearings are identical - target at infinite range.");

    if (angleDiff >= 90)
        throw std::runtime_error("Angle difference too large (target likely behind baseline).");

    // Convert degrees to radians: tan(90° - θ) = cot(θ)
    double angleRad = (90.0 - angleDiff) * M_PI / 180.0;
    double range = refDistance * std::tan(angleRad);

    return range;
}
