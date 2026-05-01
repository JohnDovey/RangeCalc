using System;

namespace RangeCalcMac
{
    class Program
    {
        private const string ProgVer = "1.0.13";
        private const string Header = "+------------------------------[{0,6}]----------------------------+";

        static void Main()
        {
            PrintHeader();

            try
            {
                int bearing1 = GetBearing("Bearing from Ref One to Target");
                int bearing2 = GetBearing("Bearing from Ref Two to Target");
                int refDistance = GetPositiveInt("Distance between Ref One and Ref Two (meters)");

                double range = CalculateRange(bearing1, bearing2, refDistance);

                Console.WriteLine();
                Console.ForegroundColor = ConsoleColor.Green;
                Console.WriteLine($"Range to Target: {range:F2} meters");
                Console.ResetColor();
            }
            catch (Exception ex) when (ex is not OperationCanceledException)
            {
                Console.ForegroundColor = ConsoleColor.Red;
                Console.WriteLine($"Error: {ex.Message}");
                Console.ResetColor();
            }

            Console.WriteLine("\nHit any key to exit...");
            Console.ReadKey(true);
        }

        private static void PrintHeader()
        {
            Console.BackgroundColor = ConsoleColor.Blue;
            Console.ForegroundColor = ConsoleColor.White;

            Console.WriteLine(Header, ProgVer);
            Console.WriteLine("| Range Calculator - by John Dovey <dovey.john@gmail.com>          |");
            Console.WriteLine("| View the code on GitHub [https://github.com/JohnDovey/RangeCalc] |");
            Console.WriteLine(Header, ProgVer);

            Console.ResetColor();
            Console.WriteLine("\nEnter the values as prompted (0 to exit)\n");
        }

        private static int GetBearing(string prompt)
        {
            while (true)
            {
                Console.Write($"{prompt}: ");
                string input = Console.ReadLine()?.Trim();

                if (string.IsNullOrEmpty(input))
                    continue;

                if (int.TryParse(input, out int bearing))
                {
                    if (bearing == 0) throw new OperationCanceledException();
                    if (bearing is < 1 or > 360)
                    {
                        Console.WriteLine("Error: Bearing must be between 1 and 360 degrees.");
                        continue;
                    }
                    Console.WriteLine($"\t({bearing}° from reference point)");
                    return bearing;
                }

                Console.WriteLine("Please enter a valid number.");
            }
        }

        private static int GetPositiveInt(string prompt)
        {
            while (true)
            {
                Console.Write($"{prompt}: ");
                string input = Console.ReadLine()?.Trim();

                if (string.IsNullOrEmpty(input))
                    continue;

                if (int.TryParse(input, out int value))
                {
                    if (value == 0) throw new OperationCanceledException();
                    if (value < 1)
                    {
                        Console.WriteLine("Error: Distance must be greater than zero.");
                        continue;
                    }
                    Console.ForegroundColor = ConsoleColor.Red;
                    Console.Write($"\t({value} meters between reference points)");
                    Console.ResetColor();
                    Console.WriteLine();
                    return value;
                }

                Console.WriteLine("Please enter a valid positive number.");
            }
        }

        /// **Summary:**

        /// Calculates range using the formula: Range = refDistance × cot(|bearing1 - bearing2|)
        /// 
        private static double CalculateRange(int bearing1, int bearing2, int refDistance)
        {
            int angleDiff = Math.Abs(bearing1 - bearing2);

            // Prevent division by zero or invalid angles
            if (angleDiff == 0)
                throw new ArithmeticException("Bearings are identical - target is at infinite range or calculation not possible.");
            if (angleDiff >= 90)
                throw new ArithmeticException("Angle difference too large for this method (target likely behind baseline).");

            double angleRad = (90.0 - angleDiff) * Math.PI / 180.0;
            double range = refDistance * Math.Tan(angleRad);   // tan(90° - θ) = cot(θ)

            return range;
        }
    }
}
