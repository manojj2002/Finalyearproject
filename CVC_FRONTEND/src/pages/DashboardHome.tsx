import { Grid } from "@/components/Grid"
import { TopBar } from "@/components/TopBar"
import { SimpleImageContainerList } from "@/components/ImageContainerList"
export const DashboardHome = () => {
  return (
    <div className="bg-white rounded-lg pb-4 shadow h-[200vh]">
      {/* <TopBar/> */}
    <Grid/> 
      <SimpleImageContainerList/>
    </div>
  )
}
