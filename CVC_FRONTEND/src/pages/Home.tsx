import Header from '../components/Header';
import Hero from '../components/Hero';
import { ReactLenis } from 'lenis/react';
import Feature from '../components/Feature';
import Process from '../components/Process';

import Footer from '../components/Footer';

export const Home = () => {
  return (
    <ReactLenis root>

    <div className='relative isolate overflow-hidden'>
      <Header/>
      <main>
        <Hero/>
        <Feature/>
        <Process/>
        
      </main>
      <Footer/>
    </div>
    </ReactLenis>
  )
}
