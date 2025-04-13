// import Sidebar from '@/components/Sidebar';
// import { DashboardHome } from '@/pages/DashboardHome';
// import { TopBar } from '@/components/TopBar';
// export const DashboardLayout = () => {
//   return (
//     <div className='text-stone-950 bg-stone-100 grid gap-4 p-4 grid-cols-[220px_1fr]'>
//       <Sidebar/>
//       <DashboardHome />
//     </div>
//   );
// };
import Sidebar from '@/components/Sidebar';
import { TopBar } from '@/components/TopBar';
import { Outlet } from 'react-router-dom'; // 🧠 import this

export const DashboardLayout = () => {
  return (
    <div className='text-stone-950 bg-stone-100 grid gap-4 p-4 grid-cols-[220px_1fr]'>
      <Sidebar />
      <div className="flex flex-col gap-4">
        <TopBar />
        <Outlet /> {/* 🔁 This is where the nested page like DashboardHome or ScanResultsPage goes */}
      </div>
    </div>
  );
};
