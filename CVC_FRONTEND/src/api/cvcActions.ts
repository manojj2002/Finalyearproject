// containerActions.ts
import axios from "axios";

const BASE_URL = "/api/dynamic-scan";
export const startFalcoScan = async () => {
  const token = localStorage.getItem("token");
  if (!token) throw new Error("No token found");

  await axios.get(`${BASE_URL}/logs`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  
};
export const stopFalcoScan = async () => {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("No token found");
  
    await axios.post("/api/dynamic-scan/stop-falco", null, {
      headers: { Authorization: `Bearer ${token}` }
    });
    
  };

  export const startTrivyScan = async (imageName: string ) => {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("No token found");
  
    await axios.post(`/api/static-scan/scan-image/${imageName}`, null, {
      headers: { Authorization: `Bearer ${token}` }
    });
    
  };

  export const startImagePull = async (imageName: string ) => {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("No token found");
  
    await axios.post(`/api/static-scan/pull-image/${imageName}`, null, {
      headers: { Authorization: `Bearer ${token}` }
    });
    
  };



  export const fetchUserScanResults = async () => {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("No token found");
  
    const res = await axios.get("/api/static-scan/getScanResults", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    console.log(res.data)
    return res.data; // array of scan result objects
  };

export const fetchFalcoAlerts = async () => {
  const token = localStorage.getItem("token");
    if (!token) throw new Error("No token found");
  
  const res = await axios.get("/api/dynamic-scan/getUserAlerts", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  console.log(res.data)
  return res.data;
};