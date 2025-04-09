import Sidebar from '@/components/Sidebar';
import { DashboardHome } from '@/pages/DashboardHome';

export const DashboardLayout = () => {
  return (
    <div className='text-stone-950 bg-stone-100 grid gap-4 p-4 grid-cols-[220px_1fr]'>
      <Sidebar/>
      <DashboardHome />
    </div>
  );
};
