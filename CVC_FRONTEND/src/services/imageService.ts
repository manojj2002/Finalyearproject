// src/api/getContainersByImage.ts

import axios from "axios";


export const getUserContainers = async () => {
  try {
    const res = await axios.get("/api/container/getContainerDetails")
    return res.data.containers; // Assumes backend returns { containers: [] }
  } catch (err) {
    console.error("❌ Failed to fetch containers:", err);
    return [];
  }
};

export const getUserImages = async () => {
  try {
    const res = await axios.get("/api/static-scan/getImageDetails");
    return res.data; // Assumes backend returns []Image
  } catch (err) {
    console.error("❌ Failed to fetch images:", err);
    return [];
  }
};