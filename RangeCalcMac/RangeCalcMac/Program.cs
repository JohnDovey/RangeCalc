using System;
using System.Linq.Expressions;

namespace RangeCalcMac
{
    class Program
    {
        static void Main(string[] args)
        {
            int bearing1;         // bearing from left ref point to target
            int bearing2;         // bearing from right ref point to target
            int refdistance;      // distance between two ref points _ use %u for unsigned number
            double rangetotarget; // Calculated range
            int tmpAngle;


            const string ProgVer = "1.0.12";
            string inputStr = "";
            const string OutTxt0 = "+------------------------------[{0,6}]----------------------------+";
            const string OutTxt1 = "| Range Calculator - by John Dovey <dovey.john@gmail.com>          |";
            const string OutTxt2 = "| View the code on GitHub [https://github.com/JohnDovey/RangeCalc] |";
            Console.BackgroundColor = ConsoleColor.Blue;
            Console.ForegroundColor = ConsoleColor.White;
            Console.WriteLine(OutTxt0, ProgVer);
            Console.WriteLine(OutTxt1);
            Console.WriteLine(OutTxt2);
            Console.WriteLine(OutTxt0, ProgVer);
            Console.ResetColor();

            Console.WriteLine(" Enter the values as prompted");

            try
            {
                // Bearing One
                Console.WriteLine("Bearing from Ref One to Target: ");
                inputStr = Console.ReadLine();
                bearing1 = Convert.ToInt16(inputStr);
                Console.WriteLine("\t({0} degrees bearing Ref 1 to target)", bearing1);

                if (bearing1 > 360)
                {
                    throw new ArithmeticException("Error! No more than 360 Degrees allowed");
                }

                if (bearing1 < 1)
                {
                    throw new ArithmeticException("Error! Bearing must be greater than zero");
                }

                // Bearing Two
                Console.WriteLine("Bearing from Ref Two to Target: ");
                inputStr = Console.ReadLine();
                bearing2 = Convert.ToInt16(inputStr);
                Console.WriteLine("\t({0} degrees bearing Ref 2 to target) ", bearing2);

                if (bearing2 > 360)
                {
                    throw new ArithmeticException("ANo more than 360 Degrees allowed");
                }
                if (bearing2 < 1)
                {
                    throw new ArithmeticException("Error! Bearing must be greater than zero");
                }

                // Seperation Distance
                Console.WriteLine("Distance between Ref One and Ref Two: ");
                inputStr = Console.ReadLine();
                refdistance = Convert.ToInt16(inputStr);
                Console.Write("\t(");
                Console.ForegroundColor = ConsoleColor.Red;
                Console.Write(refdistance);
                Console.ResetColor();

                Console.WriteLine(" meters between reference points) ");

                // d = (Tan (90 - (A -B))) x Ref
                tmpAngle = bearing1 - bearing2;
                if (tmpAngle < 0)
                {
                    tmpAngle *= (0 - 1);
                }
                rangetotarget = Math.Tan(90 - tmpAngle) * refdistance;
                Console.WriteLine("Range to Target: {0}", rangetotarget);
            }
            catch (Exception e)
            {
                Console.WriteLine(e.Message);
            }
            Console.Write(@"Hit any key (gently) to continue");
            Console.ReadKey();

        }
    }
}
