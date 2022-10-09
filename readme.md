# IIUMPassGo

A simple way to automate login to campus Wifi using 100% Go code

### New Update!

You can now login to iMaalum and get the files that you need ie. Financial Statements, Course Schedule, Exam TimeTable etc

## Some Context

I did the same app using C# and it was too big (138MB). So I've decided to make an app that's portable to a lot of platforms and easier to distribute. 

Feel free to add to this app


## Instructions

Clone the repo, and change directory to repo directory, and run the command **go build**

## Compatibility

Because I used Win32 libraries to display a dialog, a slight tweak needed for it to work on multiple platforms by code that references to Win32 libraries and replacies with an abstraction library
