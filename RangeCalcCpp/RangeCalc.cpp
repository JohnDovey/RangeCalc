#include <iostream>
#include <iomanip>
#include <string>
#include <cmath>
#include <limits>
#include <stdexcept>

const std::string ProgVer = "2.0.0";

// Portable pi (MSVC does not always define M_PI)
const double PI = 3.14159265358979323846;

void PrintHeader()
{
    std::cout << "+------------------------------[" << std::setw(6) << ProgVer << "]----------------------------+\n";
    std::cout << "|        Range Calculator - by John Dovey <dovey.john@gmail.com>     |\n";
    std::cout << "| View the code on GitHub [https://github.com/JohnDovey/RangeCalc]   |\n";
    std::cout << "+--------------------------------------------------------------------+\n\n";
}

// Smallest angle between two compass bearings, in (0, 180]
double CircularAngleDiff(double a, double b)
{
    double diff = std::abs(a - b);
    if (diff > 180.0)
        diff = 360.0 - diff;
    return diff;
}

double DegToRad(double deg)
{
    return deg * PI / 180.0;
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

            std::cout << "\t(" << value << " deg)\n";
            return value;
        }
        catch (const std::runtime_error&)
        {
            throw; // cancellation
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
        catch (const std::runtime_error&)
        {
            throw;
        }
        catch (...)
        {
            std::cout << "Invalid input. Please enter a positive number.\n\n";
        }
    }
}

struct TriangulationResult
{
    double rangeFromRef1;
    double rangeFromRef2;
    double angleAtTarget;
};

// Simple mode: range = baseline * cot(theta) = baseline * tan(90 - theta)
// where theta is the circular bearing difference (parallax angle).
// Assumes a symmetric / "perpendicular enough" geometry; approximate only.
// Valid for 0 < theta < 90 degrees.
TriangulationResult CalculateSimple(double b1, double b2, double baseline)
{
    double theta = CircularAngleDiff(b1, b2);

    if (theta == 0.0)
        throw std::runtime_error("Bearings are identical - cannot calculate range.");

    if (theta >= 90.0)
        throw std::runtime_error(
            "Angle difference too large for simple mode (need < 90 deg). "
            "Try Full Triangulation with a baseline bearing.");

    // cot(theta) = tan(90 - theta); convert degrees -> radians for std::tan
    double range = baseline * std::tan(DegToRad(90.0 - theta));

    return { range, range, theta };
}

// Full triangulation via law of sines.
// Requires baseline bearing (direction from Ref1 to Ref2) so the triangle
// angles at each reference point can be formed from absolute compass bearings.
TriangulationResult CalculateFull(double b1, double b2, double baseline, double baselineBearing)
{
    double angleAtTarget = CircularAngleDiff(b1, b2);

    if (angleAtTarget == 0.0)
        throw std::runtime_error("Bearings are identical - cannot calculate range.");

    // Interior angle at Ref1: between baseline (Ref1->Ref2) and LOB (Ref1->Target)
    double angleAtRef1 = CircularAngleDiff(baselineBearing, b1);

    // Third angle of the plane triangle
    double angleAtRef2 = 180.0 - angleAtTarget - angleAtRef1;

    if (angleAtRef1 <= 0.0 || angleAtRef2 <= 0.0)
        throw std::runtime_error(
            "Invalid geometry - target position not possible with these bearings.");

    // Law of sines: a/sin(A) = b/sin(B) = c/sin(C)
    //   range_from_ref1 / sin(angleAtRef2) = baseline / sin(angleAtTarget)
    double sinTarget = std::sin(DegToRad(angleAtTarget));
    if (std::abs(sinTarget) < 1e-12)
        throw std::runtime_error("Degenerate triangle (angle at target near 0 or 180).");

    double rangeFromRef1 = baseline * std::sin(DegToRad(angleAtRef2)) / sinTarget;
    double rangeFromRef2 = baseline * std::sin(DegToRad(angleAtRef1)) / sinTarget;

    if (rangeFromRef1 <= 0.0 || rangeFromRef2 <= 0.0)
        throw std::runtime_error("Invalid geometry - computed range is not positive.");

    return { rangeFromRef1, rangeFromRef2, angleAtTarget };
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

        std::cout << "\nChoose calculation method:\n";
        std::cout << "1. Simple (fast; assumes favourable geometry)\n";
        std::cout << "2. Full Triangulation (accurate; needs baseline bearing)\n";
        std::cout << "Enter choice (1 or 2): ";

        std::string choice;
        std::getline(std::cin, choice);

        TriangulationResult result;
        if (choice == "2")
        {
            int baselineBearing = GetBearing("Baseline Bearing (direction from Ref1 to Ref2)");
            result = CalculateFull(bearing1, bearing2, refDistance, baselineBearing);
        }
        else
        {
            result = CalculateSimple(bearing1, bearing2, refDistance);
        }

        std::cout << "\n";
        std::cout << std::fixed << std::setprecision(2);
        std::cout << "=== RESULTS ===\n";
        std::cout << "Range from Ref One : " << result.rangeFromRef1 << " meters\n";
        std::cout << "Range from Ref Two : " << result.rangeFromRef2 << " meters\n";
        std::cout << std::setprecision(1);
        std::cout << "Angle at Target    : " << result.angleAtTarget << " deg\n\n";
    }
    catch (const std::exception& e)
    {
        std::cout << "\nError: " << e.what() << std::endl;
    }

    std::cout << "Press Enter to exit...";
    std::cin.get();
    return 0;
}
