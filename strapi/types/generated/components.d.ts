import type { Attribute, Schema } from '@strapi/strapi';

export interface RouteRouteStop extends Schema.Component {
  collectionName: 'components_route_route_stops';
  info: {
    description: '';
    displayName: '\u041E\u0441\u0442\u0430\u043D\u043E\u0432\u043A\u0430 \u043C\u0430\u0440\u0448\u0440\u0443\u0442\u0430';
  };
  attributes: {
    order: Attribute.Integer;
    place: Attribute.Relation<
      'route.route-stop',
      'oneToOne',
      'api::place.place'
    >;
  };
}

declare module '@strapi/types' {
  export module Shared {
    export interface Components {
      'route.route-stop': RouteRouteStop;
    }
  }
}
