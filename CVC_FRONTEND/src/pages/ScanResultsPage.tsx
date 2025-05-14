import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import { fetchUserScanResults, fetchFalcoAlerts } from "@/api/cvcActions";

const TrivySeverityBar = ({ severityCount }: { severityCount: Record<string, number> }) => {
  const severityColors: Record<string, string> = {
    LOW: "bg-green-500",
    MEDIUM: "bg-yellow-400",
    HIGH: "bg-orange-500",
    CRITICAL: "bg-red-600",
  };

  const total = Object.values(severityCount).reduce((sum, val) => sum + val, 0);

  if (total === 0) {
    return <p className="text-gray-500 text-sm">No vulnerabilities found 🎉</p>;
  }

  const severityLevels = ["LOW", "MEDIUM", "HIGH", "CRITICAL"];

  return (
    <div className="w-full border rounded overflow-hidden shadow h-5 flex mt-2">
      {severityLevels.map((severity) => {
        const count = severityCount[severity] || 0;
        const percent = (count / total) * 100;
        return (
          <div
            key={severity}
            className={`${severityColors[severity] || "bg-gray-300"} h-full`}
            style={{ width: `${percent}%` }}
            title={`${severity}: ${count}`}
          />
        );
      })}
    </div>
  );
};

// export const ScanResultsPage = () => {
//   const [trivyResults, setTrivyResults] = useState<any[]>([]);
//   const [falcoAlerts, setFalcoAlerts] = useState<any[]>([]);

//   useEffect(() => {
//     const fetchResults = async () => {
//       try {
//         const [trivyResultsData, falcoResultsData] = await Promise.all([
//           fetchUserScanResults(),
//           fetchFalcoAlerts(),
//         ]);
//         setTrivyResults(trivyResultsData);
//         setFalcoAlerts(falcoResultsData);
//         toast.success("📊 Scan results loaded");
//       } catch (err) {
//         toast.error("❌ Failed to fetch scan results");
//       }
//     };

//     fetchResults();
//   }, []);

//   return (
//     <div>
//       <div className="mt-6 bg-white rounded shadow p-4">
//         <h2 className="text-lg font-semibold mb-2">🧪 Trivy Scan Results</h2>
//         {trivyResults.length === 0 ? (
//           <p className="text-sm text-gray-500">No scan results found</p>
//         ) : (
//           <ul className="space-y-2">
//             {trivyResults.map((r, i) => (
//               <li key={i} className="border p-3 rounded bg-gray-100 shadow-sm text-sm">
//                 <div><strong>Image:</strong> {r.ImageName}</div>
//                 <div><strong>Scan Time:</strong> {new Date(r.ScanTime).toLocaleString()}</div>
//                 <div className="mt-2">
//                   <strong>Severity Count:</strong>
//                   <ul className="ml-4 list-disc">
//                     {r.SeverityCount ? (
//                       Object.entries(r.SeverityCount).map(([severity, count]: any) => (
//                         <li key={severity}>{severity}: {count}</li>
//                       ))
//                     ) : (
//                       <li>No severity data available</li>
//                     )}
//                   </ul>
//                 </div>

//                 <div className="mt-2">
//                   <strong>Severity:</strong>
//                   <TrivySeverityBar severityCount={r.SeverityCount || {}} />
//                 </div>

//                 <div className="flex space-x-4 mt-2 text-xs">
//                   <span className="flex items-center"><span className="w-3 h-3 mr-1 bg-green-500 inline-block rounded-sm" /> Low</span>
//                   <span className="flex items-center"><span className="w-3 h-3 mr-1 bg-yellow-400 inline-block rounded-sm" /> Medium</span>
//                   <span className="flex items-center"><span className="w-3 h-3 mr-1 bg-orange-500 inline-block rounded-sm" /> High</span>
//                   <span className="flex items-center"><span className="w-3 h-3 mr-1 bg-red-600 inline-block rounded-sm" /> Critical</span>
//                 </div>
//               </li>
//             ))}
//           </ul>
//         )}
//       </div>

//       <div className="mt-6 bg-white rounded shadow p-4">
//         <h2 className="text-lg font-semibold mb-2">🛡️ Falco Security Alerts</h2>
//         {falcoAlerts.length === 0 ? (
//           <p className="text-sm text-gray-500">No Falco alerts found</p>
//         ) : (
//           <table className="w-full text-sm border">
//             <thead>
//               <tr className="bg-gray-200 text-left">
//                 <th className="p-2">Rule</th>
//                 <th className="p-2">Priority</th>
//                 <th className="p-2">Container</th>
//                 <th className="p-2">Process</th>
//                 <th className="p-2">Count</th>
//               </tr>
//             </thead>
//             <tbody>
//               {falcoAlerts.map((alert, i) => (
//                 <tr key={i} className="border-t">
//                   <td className="p-2">{alert.rule}</td>
//                   <td className="p-2">{alert.priority}</td>
//                   <td className="p-2">{alert["container.name"]}</td>
//                   <td className="p-2">{alert["proc.name"]}</td>
//                   <td className="p-2">{alert.count}</td>
//                 </tr>
//               ))}
//             </tbody>
//           </table>
//         )}
//       </div>
//     </div>
//   );
// };
export const ScanResultsPage = () => {
  const [trivyResults, setTrivyResults] = useState<any[]>([]);
  const [falcoAlerts, setFalcoAlerts] = useState<any[]>([]);

  useEffect(() => {
    const fetchResults = async () => {
      try {
        const [trivyResultsData, falcoResultsData] = await Promise.all([
          fetchUserScanResults(),
          fetchFalcoAlerts(),
        ]);
        setTrivyResults(trivyResultsData);
        setFalcoAlerts(falcoResultsData);
        toast.success("📊 Scan results loaded");
      } catch (err) {
        toast.error("❌ Failed to fetch scan results");
      }
    };

    fetchResults();
  }, []);

  return (
    <div>
      <div className="mt-6 bg-white rounded shadow p-4">
        <h2 className="text-lg font-semibold mb-2">🧪 Trivy Scan Results</h2>
        {!trivyResults || trivyResults.length === 0 ? (
          <p className="text-sm text-gray-500">No scan results found</p>
        ) : (
          <ul className="space-y-2">
            {trivyResults.map((r, i) => (
              <li key={i} className="border p-3 rounded bg-gray-100 shadow-sm text-sm">
                <div><strong>Image:</strong> {r.ImageName}</div>
                <div><strong>Scan Time:</strong> {new Date(r.ScanTime).toLocaleString()}</div>
                <div className="mt-2">
                  <strong>Severity Count:</strong>
                  <ul className="ml-4 list-disc">
                    {r.SeverityCount ? (
                      Object.entries(r.SeverityCount).map(([severity, count]: any) => (
                        <li key={severity}>{severity}: {count}</li>
                      ))
                    ) : (
                      <li>No severity data available</li>
                    )}
                  </ul>
                </div>

                <div className="mt-2">
                  <strong>Severity:</strong>
                  <TrivySeverityBar severityCount={r.SeverityCount || {}} />
                </div>

                <div className="flex space-x-4 mt-2 text-xs">
                  <span className="flex items-center"><span className="w-3 h-3 mr-1 bg-green-500 inline-block rounded-sm" /> Low</span>
                  <span className="flex items-center"><span className="w-3 h-3 mr-1 bg-yellow-400 inline-block rounded-sm" /> Medium</span>
                  <span className="flex items-center"><span className="w-3 h-3 mr-1 bg-orange-500 inline-block rounded-sm" /> High</span>
                  <span className="flex items-center"><span className="w-3 h-3 mr-1 bg-red-600 inline-block rounded-sm" /> Critical</span>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="mt-6 bg-white rounded shadow p-4">
        <h2 className="text-lg font-semibold mb-2">🛡️ Falco Security Alerts</h2>
        {!falcoAlerts || falcoAlerts.length === 0 ? (
          <p className="text-sm text-gray-500">No Falco alerts found</p>
        ) : (
          <table className="w-full text-sm border">
            <thead>
              <tr className="bg-gray-200 text-left">
                <th className="p-2">Rule</th>
                <th className="p-2">Priority</th>
                <th className="p-2">Container</th>
                <th className="p-2">Process</th>
                <th className="p-2">Count</th>
              </tr>
            </thead>
            <tbody>
              {falcoAlerts.map((alert, i) => (
                <tr key={i} className="border-t">
                  <td className="p-2">{alert.rule}</td>
                  <td className="p-2">{alert.priority}</td>
                  <td className="p-2">{alert["container.name"]}</td>
                  <td className="p-2">{alert["proc.name"]}</td>
                  <td className="p-2">{alert.count}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};
