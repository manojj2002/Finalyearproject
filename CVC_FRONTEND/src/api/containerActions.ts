// src/api/containerActions.ts
import axios from "axios";

const BASE_URL = "/api/container";

export const createContainerFromImage = async (image: string) => {
  const token = localStorage.getItem("token");

  if (!token) throw new Error("No token found");

  await axios.post(`/api/container/createContainer/${image}`, null, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
};

export const startContainer = async (name: string) => {
  const token = localStorage.getItem("token");

  if (!token) throw new Error("No token found");

  await axios.post(`${BASE_URL}/startContainer/${name}`, null, {
    headers: { Authorization: `Bearer ${token}` },
  });
};

export const stopContainer = async (name: string) => {
  const token = localStorage.getItem("token");

  if (!token) throw new Error("No token found");

  await axios.post(`${BASE_URL}/stopContainer/${name}`, null, {
    headers: { Authorization: `Bearer ${token}` },
  });
};

export const deleteContainer = async (name: string) => {
  const token = localStorage.getItem("token");

  if (!token) throw new Error("No token found");
  
  await axios.delete(`${BASE_URL}/deleteContainer/${name}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
};

