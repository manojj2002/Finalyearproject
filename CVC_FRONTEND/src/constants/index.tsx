import { MenuItem } from '@/types';

import {
  ChartArea,
  Building2,
  Component,
  Code,
  BetweenHorizonalEnd,
  BrainCircuit,
  Blocks,
  Terminal,
  Package,
  SquareMousePointer,
  ChartPie,
  Files,
  UserRoundPen,
  GitFork,
  LaptopMinimal,
  ArrowBigDownDash,
  CreditCard,
  Twitter,
  Github,
  Linkedin,
  Instagram,
  Youtube,
} from 'lucide-react';

import {
  feature1,
  feature2,
  
} from '@/assets';

// Header
export const navMenu: MenuItem[] = [
  {
    href: '/products',
    label: 'Products',
    submenu: [
      {
        href: '#',
        icon: <ChartArea />,
        label: 'User Analytics',
        desc: 'Powerful options to securely authenticate and manage',
      },
      {
        href: '#',
        icon: <Building2 />,
        label: 'B2B SaaS Suite',
        desc: 'Add-on features built specifically for B2B applications',
      },
      {
        href: '#',
        icon: <Component />,
        label: 'React Components',
        desc: 'Embeddable prebuilt UI components for quick and seamless integrations',
      },
      {
        href: '#',
        icon: <Code />,
        label: 'Next.js Analytics',
        desc: 'The fastest and most seamless authentication solution for Next.js',
      },
      {
        href: '#',
        icon: <BetweenHorizonalEnd />,
        label: 'AnalytiX Elements',
        desc: 'Unstyled UI primitives for endless customization. Powered by AnalytiX',
      },
      {
        href: '#',
        icon: <BrainCircuit />,
        label: 'Authentication for AI',
        desc: 'Authentication and abuse protection tailored to AI applications',
      },
    ],
  },
  {
    href: '/features',
    label: 'Features',
  },
  {
    href: '/docs',
    label: 'Docs',
    submenu: [
      {
        href: '#',
        icon: <Terminal />,
        label: 'Getting Started',
        desc: 'Powerful options to securely authenticate and manage',
      },
      {
        href: '#',
        icon: <Package />,
        label: 'Core Concepts',
        desc: 'Add-on features built specifically for B2B applications',
      },
      {
        href: '#',
        icon: <SquareMousePointer />,
        label: 'Customization',
        desc: 'Embeddable prebuilt UI components for quick and seamless integrations',
      },
      {
        href: '#',
        icon: <Blocks />,
        label: 'Official Plugins',
        desc: 'The fastest and most seamless authentication solution for Next.js',
      },
    ],
  },
  {
    href: '/pricing',
    label: 'Pricing',
  },
];

// Hero
export const heroData = {
  sectionSubtitle: 'Secure Containers, Ensure Compliance',
  sectionTitle: 'Next-Gen Container Security',
  decoTitle: 'analytics',
  sectionText:
    'Next-Gen Container Security: Protect your containers with real-time vulnerability insights and compliance monitoring.',
};

// Feature
export const featureData = {
  sectionSubtitle: 'Features',
  sectionTitle: 'Secure Your Containers with Advanced Tools',
  sectionText:
    "Leverage our platform's robust features to detect vulnerabilities, ensure compliance, and protect your containerized environments.",
  features: [
    {
      icon: <ChartPie size={32} />,
      iconBoxColor: 'bg-blue-600',
      title: 'Vulnerability Scanning',
      desc: 'Identify and prioritize vulnerabilities in container images with comprehensive scans, uncovering critical risks and CVEs.',
      imgSrc: feature1,
    },
    {
      icon: <Files size={32} />,
      iconBoxColor: 'bg-cyan-500',
      title: 'Real-Time Monitoring',
      desc: 'Monitor container runtime activity and detect threats instantly, ensuring proactive security for your infrastructure.',
      imgSrc: feature2,
    },
    {
      icon: <UserRoundPen size={32} />,
      iconBoxColor: 'bg-yellow-500',
      title: 'Compliance Assurance',
      desc: 'Validate compliance with standards like CIS and NIST, with automated checks to maintain regulatory adherence.',
    },
    {
      icon: <GitFork size={32} />,
      iconBoxColor: 'bg-red-500',
      title: 'Misconfiguration Detection',
      desc: 'Detect and resolve container and Kubernetes misconfigurations to prevent security gaps and optimize performance.',
    },
    {
      icon: <Blocks size={32} />,
      iconBoxColor: 'bg-purple-500',
      title: 'Seamless Integrations',
      desc: 'Integrate with CI/CD pipelines and cloud platforms to streamline security workflows and enhance scalability.',
    },
  ],
};

// Process
export const processData = {
  sectionSubtitle: 'How it works',
  sectionTitle: 'Simple Steps to Secure Your Containers',
  sectionText:
    'Get started with our platform using a streamlined process to protect your containerized environments.',
  list: [
    {
      icon: <LaptopMinimal size={32} />,
      title: 'Sign Up for an Account',
      text: 'Create your account to access our platform and begin securing your containers with powerful tools.',
    },
    {
      icon: <ArrowBigDownDash size={32} />,
      title: 'Integrate with Your Environment',
      text: 'Connect our platform to your CI/CD pipeline or cloud infrastructure to scan and monitor containers seamlessly.',
    },
    {
      icon: <CreditCard size={32} />,
      title: 'Start Scanning and Monitoring',
      text: 'Initiate vulnerability scans and real-time monitoring to detect threats and ensure compliance instantly.',
    },
  ],
};






// Footer
export const footerData = {
  links: [
    {
      title: 'Product',
      items: [
        {
          href: '#',
          label: 'Components',
        },
        {
          href: '#',
          label: 'Pricing',
        },
        {
          href: '#',
          label: 'Dashboard',
        },
        {
          href: '#',
          label: 'Feature requests',
        },
      ],
    },
    {
      title: 'Developers',
      items: [
        {
          href: '#',
          label: 'Documentation',
        },
        {
          href: '#',
          label: 'Discord server',
        },
        {
          href: '#',
          label: 'Support',
        },
        {
          href: '#',
          label: 'Glossary',
        },
        {
          href: '#',
          label: 'Changelog',
        },
      ],
    },
    {
      title: 'Company',
      items: [
        {
          href: '#',
          label: 'About',
        },
        {
          href: '#',
          label: 'Careers',
        },
        {
          href: '#',
          label: 'Blog',
        },
        {
          href: '#',
          label: 'Contact',
        },
      ],
    },
    {
      title: 'Legal',
      items: [
        {
          href: '#',
          label: 'Terms and Conditions',
        },
        {
          href: '#',
          label: 'Privacy Policy',
        },
        {
          href: '#',
          label: 'Data Processing Agreement',
        },
        {
          href: '#',
          label: 'Cookie manager',
        },
      ],
    },
  ],
  copyright: '© 2024 RAM',
  socialLinks: [
    {
      href: 'https://x.com/codewithsadee_',
      icon: <Twitter size={18} />,
    },
    {
      href: 'https://github.com/codewithsadee',
      icon: <Github size={18} />,
    },
    {
      href: 'https://www.linkedin.com/in/codewithsadee/',
      icon: <Linkedin size={18} />,
    },
    {
      href: 'https://www.instagram.com/codewithsadee',
      icon: <Instagram size={18} />,
    },
    {
      href: 'https://www.youtube.com/codewithsadee',
      icon: <Youtube size={18} />,
    },
  ],
};