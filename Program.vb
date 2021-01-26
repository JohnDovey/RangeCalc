Imports System

Module Program
    Sub Main()
		Dim bearing1 As Integer ' bearing from left ref point to target
		Dim bearing2 As Integer ' bearing from right ref point to target
		Dim refdistance As Integer ' distance between two ref points _ use %u for unsigned number
		Dim rangetotarget As Double ' Calculated range
		Dim tmpAngle As Integer

		Console.Write(vbLf & vbLf & " Enter the values as prompted" & vbCrLf)

		' Bearing One
		Console.Write("Bearing from Ref One to Target: ")
		bearing1 = Console.ReadLine
		Console.WriteLine(vbTab & "({0} degrees bearing Ref 1 to target) ", bearing1)

		If bearing1 > 360 Then
			Console.WriteLine("Error! No more than 360 Degrees allowed")
			Exit Sub
		End If

		If bearing1 < 1 Then
			Console.WriteLine("Error! Bearing must be greater than zero")
			Exit Sub
		End If
		' Bearing Two
		Console.Write("Bearing from Ref Two to Target: ")
		bearing2 = Console.ReadLine
		Console.WriteLine(vbTab & "({0} degrees bearing Ref 2 to target) ", bearing2)

		If bearing2 > 360 Then
			Console.WriteLine("Error! No more than 360 Degrees allowed")
			Exit Sub
		End If
		If bearing2 < 1 Then
			Console.WriteLine("Error! Bearing must be greater than zero")
			Exit Sub
		End If

		' Seperation Distance
		Console.Write("Distance between Ref One and Ref Two: ")
		refdistance = Console.ReadLine
		Console.WriteLine(vbTab & "({0} meters between reference points) ", refdistance)

		' d = (Tan (90 - (A -B))) x Ref
		tmpAngle = bearing1 - bearing2
		If tmpAngle < 0 Then
			tmpAngle *= (0 - 1)
		End If
		rangetotarget = (Math.Tan(90 - tmpAngle)) * refdistance
		Console.WriteLine("Range to Target: {0}", rangetotarget)

	End Sub
End Module
