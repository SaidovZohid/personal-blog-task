import { jwtDecode } from "jwt-decode";

export const isAuthenticated = () => {
  const token = localStorage.getItem("userToken") || sessionStorage.getItem("userToken");
  if (!token) return true; 

  try {
    const decodedToken = jwtDecode(token);
    const currentTime = Date.now() / 1000; 
    if (decodedToken.exp < currentTime) {
      localStorage.removeItem("userToken");
      sessionStorage.removeItem("userToken");
      return false;
    }
    return true;
  } catch (error) {
    console.error("Invalid token:", error);
    localStorage.removeItem("userToken");
    sessionStorage.removeItem("userToken");
    return false;
  }
};
