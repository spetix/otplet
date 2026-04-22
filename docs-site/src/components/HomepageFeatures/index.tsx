import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Easy to Use',
    description: (
      <>
        Store your seed locally with <code>otplet create</code> and use it with <code>otplet show</code>.
      </>
    ),

  },
  {
    title: 'Secure',
    description: (
      <>
        Your seed is encrypted using and RSA key protected with a passphrase.
      </>
    ),
  },
  {
    title: 'Open Source',
    description: (
      <>
        Otplet is open source and available on GitHub.
      </>
    ),
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
