import { useEffect, useState } from "react";
import { FaDocker } from "react-icons/fa";

import axios from "axios";
import { toast } from "react-toastify";
import {
  createContainerFromImage,
  startContainer,
  stopContainer,
  deleteContainer,
} from "@/api/containerActions";
import { startImagePull,startFalcoScan, stopFalcoScan,startTrivyScan} from "@/api/cvcActions";
 import { useDispatch, useSelector } from "react-redux";
import { RootState } from "@/redux/store";
import { toggleFalco } from "@/redux/Store/falcoSlice";





export const SimpleImageContainerList = () => {
  const [data, setData] = useState<any[]>([]);
  const [searchImageName, setSearchImageName] = useState("");
  const dispatch = useDispatch();
  const falcoRunning = useSelector((state: RootState) => state.falco.isRunning);


  const getAuthToken = () => localStorage.getItem("token");

  const toggleFalcoScan = async () => {
    try {
      if (falcoRunning) {
        await stopFalcoScan();
        toast.info("🛑 Falco scan stopped");
      } else {
        await startFalcoScan();
        toast.success("🛡️ Falco scan started");
      }
  
      dispatch(toggleFalco()); // This updates the global Redux state
    } catch (err) {
      toast.error("❌ Falco scan toggle failed");
    }
  };
  

  const handleSearchAndScan = async () => {
    if (!searchImageName.trim()) {
      toast.warning("⚠️ Please enter an image name");
      return;
    }
  
    try {
      await startImagePull(searchImageName.trim());
      toast.success(`🧪 Trivy scan triggered for ${searchImageName}`);
      setSearchImageName("");
    } catch (err) {
      toast.error(` The image ${searchImageName}: Already available`);
    }
  };
  
 
  const fetchData = async () => {
    try {
      const token = getAuthToken();
  
      const [imageRes, containerRes] = await Promise.all([
        axios.get("/api/static-scan/getImageDetails", {
          headers: { Authorization: `Bearer ${token}` },
        }),
        axios.get("/api/container/getContainerDetails", {
          headers: { Authorization: `Bearer ${token}` },
        }),
      ]);
  
      const images = imageRes.data || [];
      const containers = Array.isArray(containerRes.data?.containers)
        ? containerRes.data.containers
        : [];
  
      const grouped = images.map((img: any) => {
        const matchingContainers = containers.filter(
          (c: any) => c.Image === img.label
        );
  
        return {
          image: img,
          containers: matchingContainers,
        };
      });
  
      setData(grouped);
    } catch (err) {
      console.error("❌ Failed to fetch image/container data:", err);
    }
  };
  

  

  useEffect(() => {
    fetchData();
  }, []);

  const handleToggleStartStop = async (container: any) => {
    const name = container.Name;
    const isRunning = container.State === "running";

    try {
      if (isRunning) {
        await stopContainer(name);
        toast.info(`Stopped container ${name}`);
      } else {
        await startContainer(name);
        toast.success(`🚀 Started container ${name}`);
      }

      setTimeout(() => {
        fetchData();
      }, 1500);
    } catch (err) {
      toast.error(`❌ Failed to ${isRunning ? "stop" : "start"} container ${name}`);
    }
  };

  const handleDelete = async (name: string) => {
    try {
      await deleteContainer(name);
      toast.warning(`🗑️ Deleting container ${name}`);
      setTimeout(() => fetchData(), 1500);
    } catch (err) {
      toast.error(`❌ Failed to delete container ${name}`);
    }
  };
  const handleDownloadPDF = (imageName: string) => {
    // Implement your PDF download logic here
    toast.info(`Downloading entire PDF for ${imageName}`);
  };

  const handleDownloadSummary = (imageName: string) => {
    // Implement your summary download logic here
    toast.info(`Downloading summary for ${imageName}`);
  };
  return (
    <div className="p-4">
      <h1 className="text-xl font-bold mb-4">🐳 Images & Containers</h1>
      <div className="mb-6 flex justify-between items-center">
  <div className="flex items-center space-x-4">
    <button
      className={`px-4 py-2 rounded text-sm ${
        falcoRunning ? "bg-red-500 hover:bg-red-600" : "bg-purple-500 hover:bg-purple-600"
      } text-white`}
      onClick={toggleFalcoScan}
    >
      {falcoRunning ? "🛑 Stop Falco Scan" : "🛡️ Start Falco Scan"}
    </button>


  </div>
</div>
<div className="mb-6 flex items-center space-x-2">
  <input
    type="text"
    placeholder="Enter image name to scan..."
    className="border p-2 rounded text-sm w-72"
    value={searchImageName}
    onChange={(e) => setSearchImageName(e.target.value)}
  />
  <button
    className="bg-blue-600 hover:bg-blue-700 text-white px-3 py-2 rounded text-sm"
    onClick={handleSearchAndScan}
  >
    🔍 Pull Image
  </button>
</div>


      {data.map((group, idx) => (
        <div key={idx} className="border rounded p-4 mb-4 shadow-sm bg-gray-100">
          <h2 className="text-lg font-semibold mb-2 flex justify-between items-center">
            <span>
              📦 {group.image.label} ({group.image.repository})
            </span>
            <button
              className="bg-blue-500 text-white px-3 py-1 rounded hover:bg-blue-600 text-sm"
              onClick={() =>
                createContainerFromImage(group.image.label).then(() => {
                  toast.success("Container created");
                  setTimeout(() => fetchData(), 1000);
                })
              }
            >
              ➕ Create Container
            </button>
            <button
  className="bg-blue-500 text-white px-3 py-1 rounded hover:bg-blue-600 text-sm"
  onClick={async () => {
    try {
      await startTrivyScan(group.image.label);
      toast.success("🧪 Trivy scan triggered");
    } catch (err) {
      toast.error("Already Scanned");
    }
  }}
>
  🧪 Scan With Trivy
</button>

          </h2>

          {group.containers.length > 0 ? (
            <ul className="ml-6 list-disc text-sm">
              {group.containers.map((c: any) => (
                <li key={c.ID} className="flex items-center justify-between mb-2">
                  <div>
                    🐳 <strong>{c.Name}</strong> | ID: {c.ID?.slice(0, 12)} | Status: {c.State}
                  </div>
                  <div className="space-x-2">
                    <button
                      className={`${
                        c.State === "running"
                          ? "bg-yellow-500 hover:bg-yellow-600"
                          : "bg-green-500 hover:bg-green-600"
                      } text-white px-2 py-1 rounded`}
                      onClick={() => handleToggleStartStop(c)}
                    >
                      {c.State === "running" ? "Stop" : "Start"}
                    </button>
                    <button
                      className="bg-red-500 text-white px-2 py-1 rounded hover:bg-red-600"
                      onClick={() => handleDelete(c.Name)}
                    >
                      Delete
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-gray-500 ml-4">
              No containers found for this image
            </p>
          )}
        </div>
      ))}
    </div>
    
  );
};
