import { Grid } from "@/components/Grid"
import { TopBar } from "@/components/TopBar"


export const DashboardHome = () => {
  return (
    <div className="bg-white rounded-lg pb-4 shadow h-[200vh]">
      <TopBar/>
      <Grid/>
    </div>
  )
}
