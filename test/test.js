let schedule = require('./schedule.json');
let students = require('./students.json');

// console.table(schedule)

// Display first semester 
console.table(schedule.block1)
console.table(schedule.block2)
console.table(schedule.block3)
console.table(schedule.block4)

// Display second semester 
console.table(schedule.block5)
console.table(schedule.block6)
console.table(schedule.block7)
console.table(schedule.block8)

// Display each students schedule
// for (let i = 0; i < students.length; i++) {
//   console.table(students[i].schedule)
// }
