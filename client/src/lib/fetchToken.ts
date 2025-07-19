import axios from "axios"
import { BACKEND_URL } from "./constants"

export const fetchToken = () => {
  try {
    setInterval(async () => {
			const {
				data: { data: token },
      } = await axios.post(`${BACKEND_URL}/refresh`);
      
      console.log("token:", token)

			localStorage.setItem("token", token);
		}, 840000);
  } catch (error) {
    console.error(error)
  }
}
