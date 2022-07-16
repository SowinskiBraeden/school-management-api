#! python
import json
import random
from inspect import currentframe
from string import hexdigits

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

exists = lambda n : True if n not in ('', None) else False
getLineNumber = lambda : currentframe().f_back.f_lineno

# Takes in information to create or add a new conflict
# Returns if the particular student has a previous error
def newConflict(pupilNum: str, email: str, conflictType: str, code: str, description: str, logs: dict) -> bool:
  exists = True if pupilNum in logs else False
  log = {
    "Pupil #": pupilNum,
    "Email": email,
    "Type": conflictType,
    "Code": code,
    "Conflict": description
  }
  if exists: logs[pupilNum].append(log)
  else: logs[pupilNum] = [log]
  return exists if conflictType == "Critical" else False

minReq, median, classCap = 18, 24, 30
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

# These are the codes for Flex (spare) blocks
# Semester 1 and 2
flex = ("XAT--12A-S", "XAT--12B-S")

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
  
  def equal(l: list) -> list: # Used to equalize list of numbers
    q,r = divmod(sum(l),len(l))
    return [q+1]*r + [q]*(len(l)-r)


  # Step 1 - Calculate which classes can run
  for student in students:
    # Tally class request
    for request in (request for request in student["requests"] if not request["alt"] and request["CrsNo"] not in flex):
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
    classRunCount = activeCourses[index]["Requests"] // median
    remaining = activeCourses[index]["Requests"] % median

    # Put # of classRunCount classes in emptyClasses
    for j in range(classRunCount):
      emptyClasses[index][f"{index}-{hexdigits[j]}"] = {
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
          emptyClasses[index][f"{index}-{hexdigits[j]}"]["expectedLen"] += 1
          remaining -= 1

    # If we can create a class using remaining, but no other classes
    # exists, create class, and do not equalize
    elif remaining >= minReq:
      # Create a class using remaining
      emptyClasses[index][f"{index}-{hexdigits[classRunCount]}"] = {
        "CrsNo": index,
        "Description": activeCourses[index]["Description"],
        "expectedLen": remaining
      }

      classRunCount += 1
      if classRunCount >= 2:
        # Equalize (level) class expectedLen's
        expectedLengths = [emptyClasses[index][f"{index}-{hexdigits[j]}"]["expectedLen"] for j in range(classRunCount)]
        newExpectedLens = equal(expectedLengths)
        for j in range(len(newExpectedLens)):
          emptyClasses[index][f"{index}-{hexdigits[j]}"]["expectedLen"] = newExpectedLens[j]

    # Else if we can't fit remaining in open slots in existing classes
    # and it is unable to create its own class,
    # and requiered number to make a class is less than the max number we can provide from existing classes
    elif minReq - remaining < classRunCount * (median - minReq):
      # Take 1 from each class till min requirment met
      for j in range(classRunCount):
        emptyClasses[index][f"{index}-{hexdigits[j]}"]["expectedLen"] -= 1
        remaining += 1
        if remaining == minReq: break

      # Create a class using remaining
      emptyClasses[index][f"{index}-{hexdigits[classRunCount]}"] = {
        "CrsNo": index,
        "Description": activeCourses[index]["Description"],
        "expectedLen": remaining
      }
      
      classRunCount += 1

      # Equalize (level) class expectedLen's
      expectedLengths = [emptyClasses[index][f"{index}-{hexdigits[j]}"]["expectedLen"] for j in range(classRunCount)]
      newExpectedLens = equal(expectedLengths)
      for j in range(len(newExpectedLens)):
        emptyClasses[index][f"{index}-{hexdigits[j]}"]["expectedLen"] = newExpectedLens[j]

    else:
      # In the case that the remaining requests are unable to be resolved
      # Fill as many requests into class as possible, any left that can't fit,
      # Will need to be ignored so later we can fold them into their alternative
      # choices
      for j in range(classRunCount):
        if remaining == 0: break
        if emptyClasses[index][f"{index}-{hexdigits[j]}"]["expectedLen"] < classCap: 
          emptyClasses[index][f"{index}-{hexdigits[j]}"]["expectedLen"] += 1
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
    for request in (request for request in student["requests"] if not request["alt"] and request["CrsNo"] not in flex):
      course = request["CrsNo"]
      getAvailableCourse = True
      isAlt = False
      while getAvailableCourse:
        if course in emptyClasses:
          # if course exists, get first available class
          for cname in emptyClasses[course]:
            if cname in selectedCourses:
              if isAlt and emptyClasses[course][cname]["expectedLen"] < classCap:
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
    if stepType == 0: return 5 if offset == 0 or offset == -4 else -4
    
    # stepType 1 is for stepping between second and first semester
    elif stepType == 1: return -5 if offset == 0 or offset == 6 else 6

    # Return Error if code is altered to cause error
    else: raise SystemExit(f"Invalid 'stepType' in func 'stepIndex' line {getLineNumber()}")

  # Create copy for step 6
  courseRunInfoCopy = dict(courseRunInfo)

  while len(allClassRunCounts) > 0:
    # Get highest resource class (most times run)
    index = allClassRunCounts.index(min(allClassRunCounts))
    course = list(courseRunInfo)[index]

    # Tally first and second semester
    sem1, sem2 = 0, 0
    sem1List, sem2List = {}, {}
    for i in range(1, 11):
      if i <= 5:
        sem1 += len(running[f"block{i}"])
        sem1List[f"block{i}"] = running[f"block{i}"]
      elif i > 5:
        sem2 += len(running[f"block{i}"])
        sem2List[f"block{i}"] = running[f"block{i}"]

    # If there is more than one class Running
    if allClassRunCounts[index] > 1:
      blockIndex = 0 if sem1 <= sem2 else 5
      stepType = 0 if sem1 <= sem2 else 1
      offset = 0

      # Spread classes throughout both semesters
      for i in range(courseRunInfo[course]["Total"]):
        cname = f"{course}-{hexdigits[i]}"
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

      if sem1 <= sem2:
        for block in sem1List:
          semBlocks.append(len(block))

      # If sem1 is less than or equal to sem2, add to sem1
      if sem1 <= sem2: [semBlocks.append(len(block)) for block in sem1List]

      # If sem2 is less than sem1, add to sem2
      elif sem1 > sem2: [semBlocks.append(len(block)) for block in sem2List]

      # Get block with least classes
      leastBlock = semBlocks.index(min(semBlocks))
      cname = f"{course}-0"

      running[f"block{leastBlock+offset}"][cname] = {
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


  # Step 6 - Evaluate, move students to fix conflicts
  conflictLogs = {}
  criticalCount, acceptableCount = 0, 0
  c_mc_count, c_cr_count, a_mc_count = 0, 0, 0
  studentsCritical, studentsAcceptable = 0, 0

  for student in students:
    blocks = [student["schedule"][block] for block in student["schedule"]]
    hasConflicts = True if sum(1 for b in blocks if len(b)>1) > 0 else False

    # If there is no conflicts
    # and classes inserted to is equal to expectedClasses
    # or classes the student is inserted to is missing
    # no more than two classes:
    # continue to next student
    if not hasConflicts and student["classes"] == student["expectedClasses"]: continue
    elif not hasConflicts and (student["expectedClasses"]-2) <= student["classes"] < student["expectedClasses"]:
      a_mc_count += 1
      acceptableCount += 1
      if not newConflict(student["Pupil #"], "", "Acceptable", "A-MC", "Missing 1-2 Classses", conflictLogs): studentsAcceptable += 1
      continue
    
    studentData = {
      "Pupil #": student["Pupil #"],
      "index": student["studentIndex"]
    }

    if hasConflicts:
      # Clear student schedule
      for block in student["schedule"]:
        [running[block][cname]["students"].remove(studentData) for cname in student["schedule"][block]]
        student["schedule"][block] = []

      # Find what class in student schedule has least run time
      classes, runCounts = [], []
      for block in blocks:
        for cname in block:
          classes.append(cname[:-2])
          runCounts.append(courseRunInfoCopy[cname[:-2]]["Total"])

      # Rebuild student schedule
      availableBlocks = [f'block{i}' for i in range(1, 11)]
      while len(classes) > 0:
        index = runCounts.index(min(runCounts)) # Get class least run
        found = False

        # find slot for class
        for block in availableBlocks:
          if found: break
          for cname in running[block]:
            if cname[:-2] == classes[index] and len(running[block][cname]["students"]) < classCap:
              running[block][cname]["students"].append(studentData)
              student["schedule"][block].append(cname)
              availableBlocks.remove(block)
              found = True
              break

        if not found:
          # Determine all places class exists
          existsIn, existingClassNames = [], []
          for block in running:
            for cname in running[block]:
              if cname[:-2] == classes[index] and len(running[block][cname]["students"]) < classCap:
                existsIn.append(block)
                existingClassNames.append(cname)

          # Attempt to fix
          solution = False
          if len(existsIn) > 0:
            for i, existing in enumerate(existsIn):
              if solution: break
              classOut = student["schedule"][existing][0]
              for block in running:
                if solution: break
                if block == existing or block not in availableBlocks: continue
                for cname in running[block]:
                  if cname[:-2] == classOut[:-2] and len(running[block][cname]["students"]) < classCap:
                    # Move to existing class elsewhere
                    student["schedule"][block].append(cname)
                    running[block][cname]["students"].append(studentData)

                    # Overwrite old class
                    running[existing][student["schedule"][existing][0]]["students"].remove(studentData)
                    student["schedule"][existing][0] = existingClassNames[i]
                    running[existing][existingClassNames[i]]["students"].append(studentData)
                    
                    solution = True
                    break

          if not solution:
            # Try alternate
            alternates = [alt["CrsNo"] for alt in students[student["studentIndex"]]["remainingAlts"] if alt["CrsNo"] not in flex and alt["CrsNo"] in courseRunInfoCopy]
            if len(alternates) == 0: # If not alternates, create critical error
              c_cr_count += 1
              criticalCount += 1
              if not newConflict(student["Pupil #"], "", "Critical", "C-CR", "Couldn't Resolve", conflictLogs): studentsCritical += 1

            else:
              # Get alternate least run
              altRunCounts = [courseRunInfoCopy[alt]["Total"] for alt in alternates]
              altIndex = altRunCounts.index(min(altRunCounts))

              # add alternate
              classes.append(alternates[altIndex])
              runCounts.append(altRunCounts[altIndex])

              # Remove alternate from remaining alternates
              for remaining in students[student["studentIndex"]]["remainingAlts"]:
                if remaining["CrsNo"] == alternates[altIndex]:
                  students[student["studentIndex"]]["remainingAlts"].remove(remaining)

        # Remove class after inserted or failed to insert
        classes.remove(classes[index])
        runCounts.remove(runCounts[index])

    metSelfRequirements = True if student["classes"] == student["expectedClasses"] else False
    while not metSelfRequirements:
      
      if (student["expectedClasses"] - 2) <= student["classes"] < student["expectedClasses"]:
        a_mc_count += 1
        acceptableCount += 1
        if not newConflict(student["Pupil #"], "", "Acceptable", "A-MC", "Missing 1-2 Classses", conflictLogs): studentsAcceptable += 1
        break

      elif student["classes"] < (student["expectedClasses"] - 2):
        # Difference between classes inserted to and
        # expected classes is too great, attempt to fix
        if student["Pupil #"] in conflictLogs:
          c_mc_count += 1
          criticalCount += 1
          if not newConflict(student["Pupil #"], "", "Critical", "C-MC", "Missing too many Classses", conflictLogs): studentsCritical += 1

        break

      else:
        print(f"Fatal error ({getLineNumber()}): Impossible error")
        break

  finalConflictLogs = {
    "Conflicts": conflictLogs,
    "Critical": {
      "Total": criticalCount,
      "Students": studentsCritical,
      "Errors": [{
        "Total": c_mc_count,
        "Description": "Missing too many Classes",
        "Code": "C-MC"
      }, {
        "Total": c_cr_count,
        "Description": "Couldn't Resolve",
        "Code": "C-CR"
      }]
    },
    "Acceptable": {
      "Total": acceptableCount,
      "Students": studentsAcceptable,
      "Errors": [{
        "Total": a_mc_count,
        "Description": "Missing 1-2 Classes",
        "Code": "A-MC"
      }]
    }
  }

  # Log Conflict to records
  with open(conflictsDir, "w") as outfile:
    json.dump(finalConflictLogs, outfile, indent=2)

  for student in students:
    for block in student["schedule"]:
      if len(student["schedule"][block]) == 0:
        student["schedule"][block].append(flex[0]) if int(block[5:]) <= 5 else student["schedule"][block].append(flex[1])

  # Update Student records
  with open(studentsDir, "w") as outfile:
      json.dump(students, outfile, indent=2)
  
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