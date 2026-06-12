import event1Img from '../assets/event1.jpg'
import event2Img from '../assets/event2.jpg'
import event3Img from '../assets/event3.jpg'
import event4Img from '../assets/event4.jpg'
import event5Img from '../assets/event5.jpg'
import event7Img from '../assets/event7.jpg'
import event8Img from '../assets/event8.jpg'
import event9Img from '../assets/event9.jpg'

const eventImages = [
  event1Img,
  event2Img,
  event3Img,
  event4Img,
  event5Img,
  event7Img,
  event8Img,
  event9Img,
]

export const eventImagesArray = eventImages

export function getRandomEventImage(): string {
  return eventImages[Math.floor(Math.random() * eventImages.length)]
}
