#!/usr/bin/env python3
import json
import math
import random

# Import from custom utilities
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

minReq, median, classCap = 18, 24, 30
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
def generateScheduleV3(
  students: list, 
  courses: dict, 
  blockClassLimit: int=40, 
  studentsDir: str="../output/students.json", 
  conflictsDir: str="../output/conflicts.json"
  ) -> dict[str, dict]:
  
  def equal(l): # Used to equalize list of numbers
    q,r = divmod(sum(l),len(l))
    return [q+1]*r + [q]*(len(l)-r)


  # Step 1 - Calculate which classes can run
  for student in students:
    # Tally class request
    for request in (request for request in student["requests"] if not request["alt"] and request["CrsNo"] not in ["XAT--12A-S", "XAT--12B-S"]):
      code = request["CrsNo"]
      courses[code]["Requests"] += 1
      # Add course to active list if enough requests
      if courses[code]["Requests"] > minReq and courses[code]["CrsNo"] not in activeCourses:
        activeCourses[code] = courses[code]


  # Step 2 - Generate empty classes
  allClassRunCounts = []
  courseRunInfo = {} # Generated now, used in step 4
  emptyClasses = {} # List of all classes with how many students should be entered during step 3
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
  tempStudents = list(students)
  
  while len(tempStudents) > 0:
    student = tempStudents[random.randint(0, len(tempStudents)-1)]

    alternates = [request for request in student["requests"] if request["alt"]]
    for request in (request for request in student["requests"] if not request["alt"] and request["CrsNo"] not in ["XAT--12A-S", "XAT--12B-S"]):
      course = request["CrsNo"]
      getAvailableCourse = True
      isAlt = False
      while getAvailableCourse:
        if course in emptyClasses:
          # if course exists, get first available class
          for cname in emptyClasses[course]:
            if cname in selectedCourses:
              if isAlt and emptyClasses[course][cname]["expectedLen"] > classCap:
                emptyClasses[course][cname]["expectedLen"] += 1
              if len(selectedCourses[cname]["students"]) < emptyClasses[course][cname]["expectedLen"]:
                # Class exists with room for student
                selectedCourses[cname]["students"].append({
                  "Pupil #": student["Pupil #"],
                  "index": student["studentIndex"]
                })
                getAvailableCourse = False
                break
              elif len(selectedCourses[cname]["students"]) == emptyClasses[course][cname]["expectedLen"]:
                # If class is full, and is last class of that course
                if cname[len(cname)-1] == f"{len(emptyClasses[course])-1}":
                  if len(alternates) > 0:
                    # Use alternate
                    course = alternates[0]["CrsNo"]
                    alternates.remove(alternates[0])
                    isAlt = True
                    break
                  else:
                    # Force break loop, ignore and let an admin
                    # handle options to solve for missing class
                    getAvailableCourse = False
                    break
            elif cname not in selectedCourses:
              selectedCourses[cname] = {
                "students": [{
                  "Pupil #": student["Pupil #"],
                  "index": student["studentIndex"]
                }],
                "CrsNo": course,
                "Description": courses[course]["Description"]
              }
              getAvailableCourse = False
              break

        elif course not in emptyClasses:
          if len(alternates) > 0:
            # Use alternate
            course = alternates[0]["CrsNo"]
            alternates.remove(alternates[0])
            isAlt = True
          else:
            # Force break loop, ignore and let an admin
            # handle options to solve for missing class
            getAvailableCourse = False

    students[student["studentIndex"]]["remainingAlts"] = alternates
    tempStudents.remove(student)


  # Step 4 - Attempt to fit classes into timetable
  def stepIndex(offset: int, stepType: int) -> int:
    # stepType 0 is for stepping between first and second semester
    if stepType == 0:
      if offset == 0 or offset == -4: return 5
      else: return -4
    
    # stepType 1 is for stepping between second and first semester
    elif stepType == 1:
      if offset == 0 or offset == 6: return -5
      else: return 6

    # Return Error if code is altered to cause error
    else: raise SystemExit("Invalid 'stepType' in func 'stepIndex' line 225")

  while len(allClassRunCounts) > 0:
    # Get highest resource class (most times run)
    index = allClassRunCounts.index(min(allClassRunCounts))
    course = list(courseRunInfo)[index]

    # Tally first and second semester
    sem1, sem2 = 0, 0
    sem1List, sem2List = {}, {}
    for i in range(1, 6):
      sem1 += len(running[f"block{i}"])
      sem1List[f"block{i}"] = running[f"block{i}"]
    for i in range(5, 11):
      sem2 += len(running[f"block{i}"])
      sem2List[f"block{i}"] = running[f"block{i}"]

    # If there is more than one class Running
    if allClassRunCounts[index] > 1:
      blockIndex = 0 if sem1 <= sem2 else 5
      stepType = 0 if sem1 <= sem2 else 1
      offset = 0

      # Spread classes throughout both semesters
      for i in range(courseRunInfo[course]["Total"]):
        cname = f"{course}-{i}"
        classInserted = False
        while not classInserted:

          blockIndex += offset
          if len(running[list(running)[blockIndex]]) < blockClassLimit:
            running[list(running)[blockIndex]][cname] = {
              "CrsNo": course,
              "Description": emptyClasses[course][cname]["Description"],
              "students": selectedCourses[cname]["students"]
            }
            allClassRunCounts[index] -= 1
            classInserted = True

          offset = stepIndex(offset, stepType)

          if blockIndex >= 9:
            blockIndex = 0 if sem1 <= sem2 else 5
            offset = 0

    # If the class only runs once, place in semester with least classes
    elif allClassRunCounts[index] == 1:
      # Equally disperse into semesters classes
      semBlocks = []
      offset = 1

      # If sem1 is less than or equal to sem2, add to sem1
      if sem1 <= sem2:
        for block in sem1List:
          semBlocks.append(len(block))

      # If sem2 is less than sem1, add to sem2
      elif sem1 > sem2:
        offset = 5
        for block in sem2List:
          semBlocks.append(len(block))

      # Get block with least classes
      leastBlock = semBlocks.index(min(semBlocks))
      cname = f"{course}-0"

      running[f"block{leastBlock+offset}"][course] = {
        "CrsNo": course,
        "Description": emptyClasses[course][cname]["Description"],
        "students": selectedCourses[cname]["students"],
      }

      allClassRunCounts[index] -= 1

    # Remove course when fully inserted
    if allClassRunCounts[index] == 0:
      allClassRunCounts.remove(allClassRunCounts[index])
      courseRunInfo.pop(list(courseRunInfo)[index])


  # Step 5 - Fill student schedule
  for block in running:
    for cname in running[block]:
      for student in running[block][cname]["students"]:
        students[student["index"]]["schedule"][block].append(cname)
        students[student["index"]]["classes"] += 1


  # Step 6 - Evaluate, move students to fix
  conflicts = []

  for student in students:
    print(student["Pupil #"])
    blocks = [student["schedule"][block] for block in student["schedule"]]
    origin = list(block)
    exceptions = []
    count, hasConflict = 0, True
    initialCount = sum(1 for b in blocks if len(b)==1)
    if initialCount == student["classes"]:
      hasConflict = False
      
      if initialCount < student["expectedClasses"]:
        conflicts.append({
          "Pupil #": student["Pupil #"],
          "Email": "",
          "Conflict": "Missing classes"
        })

    while hasConflict:
      # Get clash
      if len(exceptions) > 0:
        blocks = [elem for i, elem in enumerate(blocks) if i not in exceptions]
        exceptCopy = list(exceptions)
        while len(exceptCopy) > 0:
          minIndex = exceptCopy.index(min(exceptCopy))  
          blocks.insert(minIndex, ['nil'])

      count = sum(1 for b in blocks if len(b)==1)

      if count < student["expectedClasses"]:
        if count == student["classes"]:
          conflicts.append({
            "Pupil #": student["Pupil #"],
            "Email": "",
            "Conflict": "Missing classes"
          })
          hasConflict = False
          break

        blockLens = [len(block) for block in blocks]
        index = blockLens.index(max(blockLens))

        blockOut = f"block{index+1}"
        done = False
        moveIndex = 0
        freeBlocks = [blockIndex for blockIndex in range(len(blocks)) if len(blocks[blockIndex]) == 0]
        while not done:
          classOut = blocks[index][moveIndex]
          found = False
          for block in running:
            if found: break
            if list(running).index(block) != index and list(running).index(block) in freeBlocks:
              for cname in running[block]:
                if cname[:-2] == blocks[index][moveIndex][:-2] and len(running[block][cname]["students"]) < classCap:
                  studentData = {
                    "Pupil #": student["Pupil #"],
                    "index": student["studentIndex"]
                  }
                  studentIndex = running[blockOut][classOut]["students"].index(studentData)
                  
                  # Update current blocks to work with
                  blocks[index].remove(classOut)
                  blocks[list(running).index(block)].append(cname)

                  # Update Final Records
                  del running[blockOut][classOut]["students"][studentIndex] # Remove final
                  running[block][cname]["students"].append(studentData)
                  
                  found, done = True, True
                  break

          if not found:
            if moveIndex < len(blocks[index])-1:
              moveIndex += 1
            elif moveIndex == len(blocks[index])-1:
              if index not in exceptions: exceptions.append(index)
              
            else:
              print('impossible err 3')
         
          elif found: done = True
          
          else:
            print('impossible err 2')

      elif count == student["expectedClasses"]:
        if len(exceptions) > 0:
          for i in range(len(exceptions)):
            blocks[i] = origin[i]

          conflicts.append({
            "Pupil #": student["Pupil #"],
            "Email": "",
            "Conflict": "More than one class per block"
          })

        hasConflict = False
      elif count > student["expectedClasses"]:
        conflicts.append({
          "Pupil #": student["Pupil #"],
          "Email": "",
          "Conflict": "More classes than expected"
        })
        hasConflict = False
      else:
        print('impossible err 1')

    # Update Student Schedule
    student["schedule"] = blocks

  # Update Student records
  with open(studentsDir, "w") as outfile:
      json.dump(students, outfile, indent=2)

  # Log Conflict to records
  with open(conflictsDir, "w") as outfile:
    json.dump(conflicts, outfile, indent=2)

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