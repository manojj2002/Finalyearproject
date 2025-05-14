// import { FiDollarSign, FiHome, FiLink, FiPaperclip, FiUsers } from 'react-icons/fi';
// import { IconType } from 'react-icons/lib';

// export const RouteSelect = () => {
//   return (
//     <div className='space-y-1'>
//       <Route Icon={FiHome} selected={true} title='Dashboard'/>
//       <Route Icon={FiUsers} selected={false} title='Home'/>
//       <Route Icon={FiPaperclip} selected={false} title='Results'/>
//       <Route Icon={FiLink} selected={false} title='Scan History'/>
//       {/* <Route Icon={FiDollarSign} selected={false} title='Finance'/> */}
//     </div>
//   );
// };

// const Route = ({
//   selected,
//   Icon,
//   title,
// }: {
//   selected: boolean;
//   Icon: IconType;
//   title: string;
// }) => {
//   return <button className={`flex items-center justify-start gap-2 w-full rounded px-2 py-1.5 text-sm transiton-[box-shadow,_background-color,_color] ${selected?"bg-white text-stone-950 shadow":"hover:bg-stone-200 bg-transparent text-stone-500 shadow-none"}`}>
//     <Icon className={selected?"text-violet-500":""}/>
//     <span>{title}</span>
//   </button>;
// };
import { NavLink } from "react-router-dom";
import { FiHome, FiPaperclip } from "react-icons/fi";

const routes = [
  { path: "/dashboard", title: "Dashboard", Icon: FiHome },
  { path: "/dashboard/result", title: "Results", Icon: FiPaperclip },
];

export const RouteSelect = () => {
  return (
    <div className="space-y-1">
      {routes.map(({ path, title, Icon }) => (
        <NavLink
          key={path}
          to={path}
          end // 🔥 ensures exact match for /dashboard (not active on /dashboard/result)
          className={({ isActive }) =>
            `flex items-center gap-2 w-full rounded px-2 py-1.5 text-sm transition-[box-shadow,_background-color,_color] ${
              isActive
                ? "bg-white text-stone-950 shadow"
                : "hover:bg-stone-200 bg-transparent text-stone-500 shadow-none"
            }`
          }
        >
          {({ isActive }) => (
            <>
              <Icon className={isActive ? "text-violet-500" : ""} />
              <span>{title}</span>
            </>
          )}
        </NavLink>
      ))}
    </div>
  );
};
