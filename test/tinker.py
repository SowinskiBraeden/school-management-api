#!/usr/bin/python3
import sys
import json
import math
from courses import mockCourses, activeCourses
from mockStudents import generateMockStudents, getSampleStudents
from generateCourses import getSampleCourses


'''
  Block 1-4 is first semester while
  block 5-8 is second semester
'''

'''
schedule example:
schedule: {
  "block1": "className",
  "block2": "className",
  "block3": "className",
  "block4": "className" 
  "block5": "className",
  "block6": "className",
  "block7": "className",
  "block8": "className"
}
'''

'''
running example:
running: {
  "block1": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block2": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block3": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block4": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block5": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block6": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block7": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block8": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}}
}
'''

# Error 1: No classes in schedule can fit this student
# Error 2: No more room in schedule for another class

# 400 by default
studentsNum = 340

err1, err2, = 0,0
minReq, median, classCap, blockClassLimit = 18, 24, 30, 12
mockStudents = []
running = {
  "block1": {},
  "block2": {},
  "block3": {},
  "block4": {},
  "block5": {},
  "block6": {},
  "block7": {},
  "block8": {}
}

# Currently V1 has an average of 9.4% success rate, with an avverage of 90.6% error rate
def generateScheduleV1():
  global err1, err2
  # Collect data and calculate schedules
  for student in mockStudents:
    # Tally class request
    for request in student["requests"]:
      mockCourses[request]["totalrequests"] += 1
      mockCourses[request]["studentindexes"].append(mockStudents.index(student))
      # Add course to active list if enough requests
      if mockCourses[request]["totalrequests"] > minReq and mockCourses[request]["code"] not in activeCourses: activeCourses[mockCourses[request]["code"]] = mockCourses[request]

  for student in mockStudents:
    alternateOffset = len(student["requests"])-8
    alternateIndex = 8
    for i in range(len(student["requests"])-alternateOffset): # Subtract x classes as they are alternatives
      currentCourse = student["requests"][i]
      generate = True
      while generate:
        # If class is allowed to run
        if currentCourse in activeCourses:
          blockIndex = 1
          courseNum = 1
          getFreeBlock = True
          while getFreeBlock:
            block = f"block{blockIndex}"
            cname = f"{mockCourses[currentCourse]['name']}-{courseNum}"
            if cname in running[block]:
              if student["schedule"][block] == "": # Add student to class
                if len(running[block][cname]["students"]) < classCap:
                  running[block][cname]["students"].append(student["name"])
                  student["schedule"][block] = cname
                  getFreeBlock = False
                else: # Find next available class or create new one
                  if blockIndex == 8:  # No available classes
                    if len(student["requests"]) == 8:
                      # This student never had any alternatives
                      # How to solve?

                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      generate = False
                    else:
                      if alternateIndex <= (len(student["requests"]) - 1):
                        currentCourse = student["requests"][alternateIndex]
                        alternateIndex += 1
                      else: # No more alterantive
                        # print(f"\nError 1: No more available classes for student {student['name']}")
                        err1 += 1
                        generate = False
                    getFreeBlock = False
                  else:
                    if (mockCourses[currentCourse]["teachers"] - courseNum) > 0:
                      courseNum += 1
                    else:
                      blockIndex += 1
                      courseNum = 1
              else:
                if blockIndex == 8: # No available classes
                  if len(student["requests"]) == 8:
                    # This student never had any alternatives
                    # How to solve?

                    # print(f"\nError 1: No more available classes for student {student['name']}")
                    err1 += 1
                    generate = False
                  else:
                    if alternateIndex <= (len(student["requests"]) - 1):
                      currentCourse = student["requests"][alternateIndex]
                      alternateIndex += 1
                    else: # No more alternatives
                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      err1 += 1
                      generate = False
                  getFreeBlock = False
                else:
                  blockIndex += 1
                  courseNum = 1
            else:
              if blockIndex == 8:
                # Class does not exists
                # Create new class in first available slot 
                blockNum = 1
                while True:
                  newBlock = f"block{blockNum}"
                  if cname not in running[newBlock] and len(running[newBlock]) < blockClassLimit:
                    if student["schedule"][newBlock] == "": # Add student to class
                      running[newBlock][cname] = {
                        "name": mockCourses[currentCourse]["name"],
                        "students": [student["name"]]
                      }
                      student["schedule"][newBlock] = cname
                      break
                    else:
                      if blockNum == 8:
                        # All student classes have been filled
                        break
                      else:
                        blockNum += 1
                        courseNum = 1
                  else:
                    if blockNum == 8:
                      # No room in school for more classes
                      # print(f"\nError 2: No more room in school for another class")
                      err2 += 1
                      break
                    else:
                      if (mockCourses[currentCourse]["teachers"] - courseNum) > 0:
                        courseNum += 1
                      else:
                        blockNum += 1
                        courseNum = 1
                break
              else:
                blockIndex += 1
                courseNum = 1
          break
        elif currentCourse not in activeCourses:
          if len(student["requests"]) == 8:
            # This student never had any alternatives
            # How to solve?

            # print(f"\nError 1: No more available classes for student {student['name']}")
            err1 += 1
            generate = False
          else:
            if alternateIndex <= (len(student["requests"]) - 1):
              currentCourse = student["requests"][alternateIndex]
              alternateIndex += 1
            else: # Out of alternatives
              # print(f"\nError 1: No more available classes for student {student['name']}")
              err1 += 1
              generate = False


# Currently V2 has an average of 0.35% success rate, with an avverage of 99.65% error rate
def generateScheduleV2():
  global err1, err2
  # Collect data and calculate schedules
  for student in mockStudents:
    # Tally class request
    for request in student["requests"]:
      mockCourses[request]["totalrequests"] += 1
      mockCourses[request]["studentindexes"].append(mockStudents.index(student))
      # Add course to active list if enough requests
      if mockCourses[request]["totalrequests"] > minReq and mockCourses[request]["code"] not in activeCourses: activeCourses[mockCourses[request]["code"]] = mockCourses[request]

  # calculate # of times to run class
  for i in range(len(activeCourses)):
    index = list(activeCourses)[i]
    classRunCount = math.floor(activeCourses[index]["totalrequests"] / classCap)
    # If there is minReq+ requests left, 1 more class could be run
    if (activeCourses[index]["totalrequests"] % classCap) > minReq: classRunCount += 1
    activeCourses[index]["classRunCount"] = classRunCount

    newBlockIndex = 1
    newCourseNum = 1
    while classRunCount > 0:
      block = f"block{newBlockIndex}"
      cname = f"{activeCourses[index]['name']}-{newCourseNum}"
      if cname not in running[block] and len(running[block]) < blockClassLimit:
        # Generate class and sub 1 from classRunCount
        running[block][cname] = {
          "name": activeCourses[index]["name"],
          "students": []
        }
        classRunCount -= 1
      else:
        if newBlockIndex == 8:
          if (activeCourses[index]["teachers"] - newCourseNum) > 0:
            newCourseNum += 1
            newBlockIndex = 1
          else:  
            # No room in school for more classes
            # print(f"\nError 2: No more room in school for another class")
            err2 += 1
            classRunCount = 0
        else: newBlockIndex += 1

  for student in mockStudents:
    alternateOffset = len(student["requests"])-8
    alternateIndex = 8
    for i in range(len(student["requests"])-alternateOffset): # Subtract x classes as they are alternatives
      currentCourse = student["requests"][i]
      generate = True
      courseNum = 1
      while generate:
        # If class is allowed to run
        if currentCourse in activeCourses:
          blockIndex = 1
          getFreeBlock = True
          cname = f"{mockCourses[currentCourse]['name']}-{courseNum}"
          while getFreeBlock:
            block = f"block{blockIndex}"
            if cname in running[block]:
              if student["schedule"][block] == "": # Add student to class
                if len(running[block][cname]["students"]) < classCap:
                  running[block][cname]["students"].append(student["name"])
                  student["schedule"][block] = cname
                  getFreeBlock = False
                else: # Find next available class or create new one
                  if blockIndex == 8:  # No available classes
                    if (activeCourses[index]["teachers"] - newCourseNum) > 0:
                      newCourseNum += 1
                      newBlockIndex = 1
                    else:
                      # No room in school for more classes
                      # print(f"\nError 2: No more room in school for another class")
                      err2 += 1
                      classRunCount = 0
                    if len(student["requests"]) == 8:
                      # This student never had any alternatives
                      # How to solve?

                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      err1 += 1
                      generate = False
                    else:
                      if alternateIndex <= (len(student["requests"]) - 1):
                        currentCourse = student["requests"][alternateIndex]
                        alternateIndex += 1
                      else: # No more alterantive
                        # print(f"\nError 1: No more available classes for student {student['name']}")
                        err1 += 1
                        generate = False
                    getFreeBlock = False
                  else: blockIndex += 1
              else:
                if blockIndex == 8: # No available classes
                  if len(student["requests"]) == 8:
                    # This student never had any alternatives
                    # How to solve?

                    # print(f"\nError 1: No more available classes for student {student['name']}")
                    err1 += 1
                    generate = False
                  else:
                    if alternateIndex <= (len(student["requests"]) - 1):
                      currentCourse = student["requests"][alternateIndex]
                      alternateIndex += 1
                    else: # No more alternatives
                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      err1 += 1
                      generate = False
                  getFreeBlock = False
                else:
                  if (mockCourses[currentCourse]["teachers"] - courseNum) > 0:
                    courseNum += 1
                  else:
                    blockIndex += 1
                    courseNum = 1
            else:
              if blockIndex == 8:
                # Class does not exists
                # Resort to alternative
                if alternateIndex <= (len(student["requests"]) - 1):
                  currentCourse = student["requests"][alternateIndex]
                  alternateIndex += 1
                else: # Out of alternatives
                  # print(f"\nError 1: No more available classes for student {student['name']}")
                  err1 += 1
                  generate = False
                break
              else:
                if (mockCourses[currentCourse]["teachers"] - courseNum) > 0:
                  courseNum += 1
                else:
                  blockIndex += 1
                  courseNum = 1
          break
        elif currentCourse not in activeCourses:
          if len(student["requests"]) == 8:
            # This student never had any alternatives
            # How to solve?

            # print(f"\nError 1: No more available classes for student {student['name']}")
            err1 += 1
            generate = False
          else:
            if alternateIndex <= (len(student["requests"]) - 1):
              currentCourse = student["requests"][alternateIndex]
              alternateIndex += 1
            else: # Out of alternatives
              # print(f"\nError 1: No more available classes for student {student['name']}")
              err1 += 1
              generate = False


# V3 differs a lot by V1/2 as it does not focus on fitting the classes
# into the time table first.
# It starts by trying to get all classes full and give all students a full class list.
# Then it starts to attempt to fit all classes into a timetable, making corretions along
# the way. Corrections being moving a students class
def generateScheduleV3(students, courses):
  # Step 1 - Calculate which classes can run
  global err1, err2
  for student in students:
    # Tally class request
    for request in student["requests"]:
      if not bool([i for i in ["XAT--12A-S", "XAT--12B-S"] if (i in request["CrsNo"])]): # Filters any requested study blocks (flex: no class block)
        code = request["CrsNo"]
        courses[code]["Requests"] += 1
        # Add course to active list if enough requests
        if courses[code]["Requests"] > minReq and courses[code]["CrsNo"] not in activeCourses:
          activeCourses[code] = courses[code]
          
  # Step 2 - Generate class list without timetable
  selectedCourses = {}
  emptyClasses = {} # List of all classes with how many students should be entered during generation
  # calculate # of times to run class
  for i in range(len(activeCourses)):
    index = list(activeCourses)[i]
    if index not in emptyClasses: emptyClasses[index] = {}
    classRunCount = math.floor(activeCourses[index]["Requests"] / median)
    remaining = activeCourses[index]["Requests"] % median

    # Put # of classRunCount classes in emptyClasses
    for j in classRunCount:
      emptyClasses[index][f"{activeCourses[index]['Description']}-{j}"] = {
        "CrsNo": index,
        "Description": activeCourses[index]["Description"],
        "expectedLen": median # Number of students expected in this class / may be altered
      }
    SomeCondition = True # Just so I don't get an error

    # If remaining fit in open slots in existing classes
    # equally disperse remaining into existing classes
    if remaining <= classRunCount * (classCap - median): 
      for j in classRunCount:
        # TODO: Equally disperse remaining into existing classes
        pass

    # Else if we can't fit remaining in, 
    # and there is more than 12, equally drain expectedLen from each classRunCount
    # to fill another class, add 1 to classRunCount        
    elif remaining > classRunCount * (classCap - median) and SomeCondition:
      pass

    # else fold students into another class (alternate?)

    #=================================================================
 
    classNum = classRunCount
    
    for student in students:
      for request in (request for request in student["requests"] if not request["alt"]):
        currentCourse = request["CrsNo"]
        if currentCourse == activeCourses[index]["CrsNo"]:
          cname = f"{activeCourses[index]['Description']}-{classNum-(classRunCount-1)}"
          if cname in selectedCourses:
            if student["Pupil #"] not in selectedCourses[cname]["students"]:
              if len(selectedCourses[cname]["students"]) < classCap:
                # Class exists and there is room
                selectedCourses[cname]["students"].append(student["Pupil #"])
              elif len(selectedCourses[cname]["students"]) == classCap:
                classRunCount -= 1
          elif cname not in selectedCourses:
            selectedCourses[cname] = {
              "students": [student["Pupil #"]],
              "CrsNo": currentCourse,
              "Description": courses[currentCourse]["Description"]
            }

  # for course in selectedCourses:
  #   print("Class: ", selectedCourses[course]["Description"])
  #   print("Students: ", len(selectedCourses[course]["students"]), "\n")
  # print("==============================")
  # print("Total: ", len(selectedCourses))

  with open("classes.json", "w") as outfile:
    json.dump(selectedCourses, outfile, indent=2)

  # Step 3 - Attempt to fit classes into timetable

  # Step 4 - Evaluate, move classes or students to fix
  return []

if __name__ == '__main__':
  if len(sys.argv) == 1:
    print("Missing argument")
    exit()
  if sys.argv[1].upper() == 'V1':
    if len(sys.argv) == 3:
      try:
        studentsNum = int(sys.argv[2])
      except:
        print("Error parsing number of students")
        exit()
    print("Processing...")
    mockStudents = generateMockStudents(studentsNum)
    generateScheduleV1()
  elif sys.argv[1].upper() == 'V2':
    if len(sys.argv) == 3:
      try:
        studentsNum = int(sys.argv[2])
      except:
        print("Error parsing number of students")
        exit()
    print("Processing...")
    mockStudents = generateMockStudents(studentsNum)
    generateScheduleV2()
  elif sys.argv[1].upper() == 'V3':
    print("Processing...")
    sampleStudents = getSampleStudents(True)
    samplemockCourses = getSampleCourses(True) 
    coursesList = generateScheduleV3(sampleStudents, samplemockCourses)
  else:
    print("Invalid argument")
    exit()

  # print("\n")
  # print(f"Error 1: x{err1}")
  # print(f"Error 2: x{err2}")
  # print("\n")

  ### Displays the number of students in each class
  # for block in running:
  #   print(f"\nBlock: {block}")
  #   print("=====================")
  #   for cl in running[block]:
  #     name = running[block][cl]["name"]
  #     students = len(running[block][cl]["students"])
  #     print(f"Class: {name} | Students: {students}")

  # Count errors in students schedules
  # errors = 0
  # for i in range(len(mockStudents)):
  #   count = 0
  #   for course in mockStudents[i]["schedule"]:
  #     if mockStudents[i]["schedule"][course]=="": count+=1
  #   if count > 0: errors += 1
  
  # print(f"{errors}/{studentsNum} student(s) have a issue with their schedule")


  # with open("schedule.json", "w") as outfile:
  #   json.dump(running, outfile, indent=2)

  # with open("students.json", "w") as outfile:
  #   json.dump(mockStudents, outfile, indent=2)

  print("Done")
