
Imports System

Module Program
	Sub Main()
		Dim Bearing1 As Integer = 0 ' bearing from left ref point to target
		Dim Bearing2 As Integer = 0 ' bearing from right ref point to target
		Dim RefDistance As Integer = 0 ' distance between two ref points
		Dim RangeToTarget As Double = 0 ' Calculated range
		Dim tmpAngle As Integer = 0
		Dim OutTxt(2) As String
		Dim ProgVer As String = "1.0.12"
		Dim inputStr As String = ""
		OutTxt(0) = "+------------------------------[{0,6}]----------------------------+"
		OutTxt(1) = "| Range Calculator - by John Dovey <dovey.john@gmail.com>          |"
		OutTxt(2) = "| View the code on GitHub [https://github.com/JohnDovey/RangeCalc] |"
		Console.WriteLine(OutTxt(0), ProgVer)
		Console.WriteLine(OutTxt(1))
		Console.WriteLine(OutTxt(2))
		Console.WriteLine(OutTxt(0), ProgVer)


		Console.WriteLine(vbLf & vbLf & " Enter the values as prompted")

		' Bearing One
		Console.Write("Bearing from Ref One to Target: ")
		inputStr = Console.ReadLine
		Bearing1 = Val(inputStr)
		Console.WriteLine(vbTab & "({0} degrees bearing Ref 1 to target) ", Bearing1)

		If Bearing1 > 360 Then
			Console.WriteLine("Error! No more than 360 Degrees allowed")
			Exit Sub
		End If

		If Bearing1 < 1 Then
			Console.WriteLine("Error! Bearing must be greater than zero")
			Exit Sub
		End If
		' Bearing Two
		Console.Write("Bearing from Ref Two to Target: ")
		inputStr = Console.ReadLine
		Bearing2 = Val(inputStr)

		Console.WriteLine(vbTab & "({0} degrees bearing Ref 2 to target) ", Bearing2)

		If Bearing2 > 360 Then
			Console.WriteLine("Error! No more than 360 Degrees allowed")
			Exit Sub
		End If
		If Bearing2 < 1 Then
			Console.WriteLine("Error! Bearing must be greater than zero")
			Exit Sub
		End If

		' Seperation Distance
		Console.Write("Distance between Ref One and Ref Two: ")
		inputStr = Console.ReadLine
		RefDistance = Val(inputStr)
		Console.WriteLine(vbTab & "({0} meters between reference points) ", RefDistance)

		' d = (Tan (90 - (A -B))) x Ref
		tmpAngle = Bearing1 - Bearing2
		If tmpAngle < 0 Then
			tmpAngle *= (0 - 1)
		End If
		RangeToTarget = (Math.Tan(90 - tmpAngle)) * RefDistance
		Console.WriteLine("Range to Target: {0}", RangeToTarget)
		Console.Write("Hit a key (softly) to continue")
		Console.ReadKey()

	End Sub
End Module
