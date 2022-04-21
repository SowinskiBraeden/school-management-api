#!/usr/bin/python3
import json
import math
import random
from util.mockStudents import getSampleStudents
from util.generateCourses import getSampleCourses


'''
  Block 1-5 is first semester while
  block 6-10 is second semester


  schedule example:
  schedule: {
    "block1": "className",
    "block2": "className",
    "block3": "className",
    ...
  }


  running example:
  running: {
    "block1": {
      classCode: {
        "className": name,
        "students": [student Name]
      },
      classCode: {
        "className": name,
        "students": [student Name]
      }
    },
    ...
  }
'''

minReq, median, classCap, blockClassLimit = 18, 24, 30, 26
mockStudents = []
activeCourses = {}
running = {
  "block1": {},
  "block2": {},
  "block3": {},
  "block4": {},
  "block5": {},
  "block6": {},
  "block7": {},
  "block8": {},
  "block9": {},
  "block10": {}
}


# V3 differs a lot by V1/2 as it does not focus on fitting the classes
# into the time table first.
# It starts by trying to get all classes full and give all students a full class list.
# Then it starts to attempt to fit all classes into a timetable, making corretions along
# the way. Corrections being moving a students class
def generateScheduleV3(students: list, courses: dict) -> dict[str, dict]:
  def equal(l): # Used to equalize list of numbers
    q,r = divmod(sum(l),len(l))
    return [q+1]*r + [q]*(len(l)-r)


  # Step 1 - Calculate which classes can run
  for student in students:
    # Tally class request
    for request in student["requests"]:
      if not bool([i for i in ["XAT--12A-S", "XAT--12B-S"] if (i in request["CrsNo"])]): # Filters any requested study blocks (flex: no class block)
        code = request["CrsNo"]
        courses[code]["Requests"] += 1
        # Add course to active list if enough requests
        if courses[code]["Requests"] > minReq and courses[code]["CrsNo"] not in activeCourses:
          activeCourses[code] = courses[code]


  # Step 2 - Generate empty classes
  allClassRunCounts = []
  courseRunInfo = {} # Generated now, used in step 4
  emptyClasses = {} # List of all classes with how many students should be entered during generation
  # calculate # of times to run class
  for i in range(len(activeCourses)):
    index = list(activeCourses)[i]
    if index not in emptyClasses: emptyClasses[index] = {}
    classRunCount = math.floor(activeCourses[index]["Requests"] / median)
    remaining = activeCourses[index]["Requests"] % median

    # Put # of classRunCount classes in emptyClasses
    for j in range(classRunCount):
      emptyClasses[index][f"{index}-{j}"] = {
        "CrsNo": index,
        "Description": activeCourses[index]["Description"],
        "expectedLen": median # Number of students expected in this class / may be altered
      }

    # If remaining fit in open slots in existing classes
    if remaining <= classRunCount * (classCap - median):
      # Equally disperse remaining into existing classes
      while remaining > 0:
        for j in range(classRunCount):
          if remaining == 0: break
          emptyClasses[index][f"{index}-{j}"]["expectedLen"] += 1
          remaining -= 1

    # If we can create a class using remaining, but no other classes
    # exists, create class, and do not equalize
    elif remaining >= minReq:
      # Create a class using remaining
      emptyClasses[index][f"{index}-{classRunCount}"] = {
        "CrsNo": index,
        "Description": activeCourses[index]["Description"],
        "expectedLen": remaining
      }

      classRunCount += 1
      if classRunCount >= 2:
        # Equalize (level) class expectedLen's
        expectedLengths = [emptyClasses[index][f"{index}-{j}"]["expectedLen"] for j in range(classRunCount)]
        newExpectedLens = equal(expectedLengths)
        for j in range(len(newExpectedLens)):
          emptyClasses[index][f"{index}-{j}"]["expectedLen"] = newExpectedLens[j]

    # Else if we can't fit remaining in open slots in existing classes
    # and it is unable to create its own class,
    # and requiered number to make a class is less than the max number we can provide from existing classes
    elif minReq - remaining < classRunCount * (median - minReq):
      # Take 1 from each class till min requirment met
      for j in range(classRunCount):
        emptyClasses[index][f"{index}-{j}"]["expectedLen"] -= 1
        remaining += 1
        if remaining == minReq: break

      # Create a class using remaining
      emptyClasses[index][f"{index}-{classRunCount}"] = {
        "CrsNo": index,
        "Description": activeCourses[index]["Description"],
        "expectedLen": remaining
      }
      
      classRunCount += 1

      # Equalize (level) class expectedLen's
      expectedLengths = [emptyClasses[index][f"{index}-{j}"]["expectedLen"] for j in range(classRunCount)]
      newExpectedLens = equal(expectedLengths)
      for j in range(len(newExpectedLens)):
        emptyClasses[index][f"{index}-{j}"]["expectedLen"] = newExpectedLens[j]

    else:
      # In the case that the remaining requests are unable to be resolved
      # Fill as many requests into class as possible, any left that can't fit,
      # Will need to be ignored so later we can fold them into their alternative
      # choices
      for j in range(classRunCount):
        if emptyClasses[index][f"{index}-{j}"]["expectedLen"] < classCap and remaining > 0: 
          emptyClasses[index][f"{index}-{j}"]["expectedLen"] += 1
          remaining -= 1

    courseRunInfo[index] = {
      "Total": classRunCount,
      "CrsNo": index
    }
    allClassRunCounts.append(classRunCount)

  # Step 3 Fill emptyClasses with Students
  selectedCourses = {}
  tempStudents = students
  
  while len(tempStudents) > 0:
    student = tempStudents[random.randint(0, len(students)-1)]

    alternates = [request for request in student["requests"] if request["alt"]]
    altOffset = None
    if len(alternates) > 0: altOffset = 0
    for request in (request for request in student["requests"] if not request["alt"] and request not in ["XAT--12A-S", "XAT--12B-S"]):
      course = request["CrsNo"]
      getAvailableCourse = True
      while getAvailableCourse:
        if course in emptyClasses: 
          # if course exists, get first available class
          for cname in emptyClasses[course]:
            if cname in selectedCourses:
              if len(selectedCourses[cname]["students"]) < emptyClasses[course][cname]["expectedLen"]:
                # Class exists with room for student
                selectedCourses[cname]["students"].append(student["Pupil #"])
                getAvailableCourse = False
                break
              elif len(selectedCourses[cname]["students"]) == emptyClasses[course][cname]["expectedLen"]:
                # If class is full, and is last class of that course
                if cname[len(cname)-1] == f"{len(emptyClasses[course])-1}":
                  if altOffset is not None and altOffset <= len(alternates)-1:
                    # Use alternate
                    course = alternates[altOffset]["CrsNo"]
                    altOffset += 1
                    break
                  else:
                    # Force break loop, ignore and let an admin
                    # handle options to solve for missing class
                    getAvailableCourse = False
                    break

            elif cname not in selectedCourses:
              selectedCourses[cname] = {
                "students": [student["Pupil #"]],
                "CrsNo": course,
                "Description": courses[course]["Description"]
              }
              getAvailableCourse = False
              break

        elif course not in emptyClasses:
          if altOffset is not None and altOffset <= len(alternates)-1:
            # Use alternate
            course = alternates[altOffset]["CrsNo"]
            altOffset += 1
          else:
            # Force break loop, ignore and let an admin
            # handle options to solve for missing class
            getAvailableCourse = False
          
    tempStudents.remove(student)

  # Step 4 - Attempt to fit classes into timetable
  while len(allClassRunCounts) > 0:
    # Get lowest resource class (least times run)
    index = allClassRunCounts.index(min(allClassRunCounts))
    course = list(courseRunInfo)[index]

    # Spread classes throughout both semesters
    blockIndex = 0
    offset = 0
    for i in range(courseRunInfo[list(courseRunInfo)[index]]["Total"]):
      cname = f"{course}-{i}"  
      blockIndex += offset
      running[list(running)[blockIndex]][cname] = {
        "CrsNo": cname[:len(cname)-2],
        "Description": course,
        "students": selectedCourses[cname]["students"]
      }
      if offset == 0 or offset == -4: offset = 5
      else: offset == -4

    # Remove course when fully inserted
    if allClassRunCounts[index] == 0:
      allClassRunCounts.remove(allClassRunCounts[index])
      courseRunInfo.pop(list(courseRunInfo)[index])

  # This will require knowing how many classrooms there are,
  # how many classrooms for a type of class, computer, wood,
  # metal, art, drama, music.
  # Languages, socials, math, science, and English all share 
  # the same classroom type and can be interchanged.
  
  # Teachers can be figured out at the end of this step,
  # The school can hire more teachers in the area required,
  # so we do not need to worry about how many tpye of classes
  # are in the same block.
  # For example: we don't need to worry if there is 4 english
  # classes in the same block. even though the school may only
  # have 2 english teachers. 


  # Step 5 - Evaluate, move classes or students to fix
  return running

if __name__ == '__main__':
  print("Processing...")

  sampleStudents = getSampleStudents(True)
  samplemockCourses = getSampleCourses(True) 
  timetable = {}
  timetable["Version"] = 3
  timetable["timetable"] = generateScheduleV3(sampleStudents, samplemockCourses)

  with open("../output/timetable.json", "w") as outfile:
    json.dump(timetable, outfile, indent=2)

  print("Done")