import React, { Suspense, lazy } from 'react';
import { RouteObject, createBrowserRouter, redirect } from 'react-router-dom';

import Layout from '@/pages/Layout';
import ErrorBoundary from '@/pages/50X';
import baseRoutes, { RouteNode } from '@/router/routes';
import { floppyNavigation } from '@/utils';

const routes: RouteObject[] = [];

const routeWrapper = (routeNodes: RouteNode[], root: RouteObject[]) => {
  routeNodes.forEach((rn) => {
    if (rn.page === 'pages/Layout') {
      rn.element = <Layout />;
      rn.errorElement = <ErrorBoundary />;
    } else {
      /**
       * cannot use a fully dynamic import statement
       * ref: https://webpack.js.org/api/module-methods/#import-1
       */
      rn.page = rn.page.replace('pages/', '');
      const Ctrl = lazy(() => import(`@/pages/${rn.page}`));
      rn.element = (
        <Suspense>
          <Ctrl />
        </Suspense>
      );
    }
    root.push(rn);
    if (rn.guard) {
      const refLoader = rn.loader;
      const refGuard = rn.guard;
      rn.loader = async (args) => {
        const gr = await refGuard(args);
        if (gr?.redirect && floppyNavigation.differentCurrent(gr.redirect)) {
          return redirect(gr.redirect);
        }

        let lr;
        if (typeof refLoader === 'function') {
          lr = await refLoader(args);
        }
        return lr;
      };
    }
    const children = Array.isArray(rn.children) ? rn.children : null;
    if (children) {
      rn.children = [];
      routeWrapper(children, rn.children);
    }
  });
};

routeWrapper(baseRoutes, routes);

export { routes, createBrowserRouter };
