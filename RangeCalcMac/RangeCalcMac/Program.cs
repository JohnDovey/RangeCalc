using System;

namespace RangeCalcMac
{
    class Program
    {
        private const string ProgVer = "2.0.1 - Full Triangulation";

        static void Main()
        {
            PrintHeader();

            try
            {
                int bearing1 = GetBearing("Bearing from Ref One to Target");
                int bearing2 = GetBearing("Bearing from Ref Two to Target");
                int baselineDistance = GetPositiveInt("Baseline Distance (Ref1 to Ref2 in meters)");

                Console.WriteLine("\nChoose calculation method:");
                Console.WriteLine("1. Simple (fast, assumes perpendicular baseline)");
                Console.WriteLine("2. Full Triangulation (more accurate + Baseline Bearing)");
                Console.Write("Enter choice (1 or 2): ");
                string choice = Console.ReadLine()?.Trim();

                TriangulationResult result;

                if (choice == "2")
                {
                    int baselineBearing = GetBearing("Baseline Bearing (direction from Ref1 to Ref2)");
                    result = CalculateFullTriangulation(bearing1, bearing2, baselineDistance, baselineBearing);
                }
                else
                {
                    result = CalculateSimple(bearing1, bearing2, baselineDistance);
                }

                Console.WriteLine();
                Console.ForegroundColor = ConsoleColor.Green;
                Console.WriteLine("=== RESULTS ===");
                Console.WriteLine($"Range from Ref One : {result.RangeFromRef1:F2} meters");
                Console.WriteLine($"Range from Ref Two : {result.RangeFromRef2:F2} meters");
                Console.WriteLine($"Angle at Target    : {result.AngleAtTarget:F1}°");
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
            Console.WriteLine("+------------------------------[{0,6}]----------------------------+", ProgVer);
            Console.WriteLine("|        Range Calculator - Full Triangulation Edition           |");
            Console.WriteLine("| View the code on GitHub [https://github.com/JohnDovey/RangeCalc] |");
            Console.WriteLine("+--------------------------------------------------------------------+");
            Console.ResetColor();
            Console.WriteLine("\nEnter the values as prompted (0 to exit)\n");
        }

        // GetBearing and GetPositiveInt methods (unchanged from your version)
        private static int GetBearing(string prompt)
        {
            while (true)
            {
                Console.Write($"{prompt}: ");
                string input = Console.ReadLine()?.Trim();

                if (string.IsNullOrEmpty(input)) continue;

                if (int.TryParse(input, out int bearing))
                {
                    if (bearing == 0) throw new OperationCanceledException();
                    if (bearing is < 1 or > 360)
                    {
                        Console.WriteLine("Error: Bearing must be between 1 and 360 degrees.");
                        continue;
                    }
                    Console.WriteLine($"\t({bearing}°)");
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

                if (string.IsNullOrEmpty(input)) continue;

                if (int.TryParse(input, out int value))
                {
                    if (value == 0) throw new OperationCanceledException();
                    if (value < 1)
                    {
                        Console.WriteLine("Error: Distance must be greater than zero.");
                        continue;
                    }
                    Console.ForegroundColor = ConsoleColor.Red;
                    Console.Write($"\t({value} meters)");
                    Console.ResetColor();
                    Console.WriteLine();
                    return value;
                }
                Console.WriteLine("Please enter a valid positive number.");
            }
        }

        private static TriangulationResult CalculateSimple(int b1, int b2, int baseline)
        {
            int diff = Math.Abs(b1 - b2);
            if (diff == 0) throw new ArithmeticException("Bearings identical.");
            if (diff >= 90) Console.WriteLine("Warning: Large angle - accuracy reduced.");

            double rad = (90.0 - diff) * Math.PI / 180.0;
            double range = baseline * Math.Tan(rad);

            return new TriangulationResult
            {
                RangeFromRef1 = range,
                RangeFromRef2 = range, // Approximate
                AngleAtTarget = diff
            };
        }

        private static TriangulationResult CalculateFullTriangulation(int bearing1, int bearing2, 
            int baseline, int baselineBearing)
        {
            // Angle at target (parallax)
            double angleAtTarget = Math.Abs(bearing1 - bearing2);
            if (angleAtTarget > 180) angleAtTarget = 360 - angleAtTarget;

            if (angleAtTarget == 0)
                throw new ArithmeticException("Bearings are identical.");

            // Angle at Ref1 = difference between baseline bearing and bearing1
            double angleAtRef1 = Math.Abs(baselineBearing - bearing1);
            if (angleAtRef1 > 180) angleAtRef1 = 360 - angleAtRef1;

            double angleAtRef2 = 180.0 - angleAtTarget - angleAtRef1;

            if (angleAtRef2 <= 0)
                throw new ArithmeticException("Invalid geometry - target position not possible.");

            // Law of Sines
            double rangeFromRef1 = baseline * Math.Sin(angleAtRef2 * Math.PI / 180.0) / Math.Sin(angleAtTarget * Math.PI / 180.0);
            double rangeFromRef2 = baseline * Math.Sin(angleAtRef1 * Math.PI / 180.0) / Math.Sin(angleAtTarget * Math.PI / 180.0);

            return new TriangulationResult
            {
                RangeFromRef1 = rangeFromRef1,
                RangeFromRef2 = rangeFromRef2,
                AngleAtTarget = angleAtTarget
            };
        }
    }

    public class TriangulationResult
    {
        public double RangeFromRef1 { get; set; }
        public double RangeFromRef2 { get; set; }
        public double AngleAtTarget { get; set; }
    }
}
