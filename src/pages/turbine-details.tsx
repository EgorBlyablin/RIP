import { useEffect, type FC } from "react";
import { Col, Container, Image, Row } from "react-bootstrap";
import { useParams } from "react-router";

import { Breadcrumbs } from "../components/breadcrumbs";
import { ROUTE_LABELS, ROUTES } from "../routes";
import { useDispatch, useSelector } from "../store";

import defaultImage from "../assets/default.jpg";
import { getTurbinesList } from "../slices/turbines";
import { BeatLoader } from "react-spinners";


export const TurbineDetailsPage: FC = () => {
    const dispatch = useDispatch();
    const { turbineId } = useParams()
    const { turbines, loading } = useSelector((state) => state.turbines);

    useEffect(() => {
        dispatch(getTurbinesList())
    }, [])

    const turbine = turbines.find((turbine) => turbine.id == turbineId)

    return <Container>
        {
            !loading
                ? <>
                    <Breadcrumbs breadcrumbs={[
                        { label: ROUTE_LABELS.HOME, path: ROUTES.HOME },
                        { label: ROUTE_LABELS.TURBINES_LIST, path: ROUTES.TURBINES_LIST },
                        { label: turbine?.title ?? "", path: "" }
                    ]} />
                    <Row style={{ backgroundColor: "#F5F7F7", border: "none", borderRadius: "6px", overflow: "hidden" }}>
                        <Col xxl={3} lg={4} md={6} xs={12} className="p-0 m-0">
                            <Row>
                                <Image fluid src={turbine?.Image || defaultImage} alt={`Изображение ${turbine?.title}`} />
                                <div style={{ padding: 36 }}>
                                    <div className="d-flex align-items-center mb-2">
                                        <span className="me-2 fw-bold">Максимальная мощность:</span>
                                        <span className="flex-grow-1 border-bottom  border-dotted"></span>
                                        <span className="ms-2">{(turbine?.power as number / 1000).toPrecision(2)} кВт</span>
                                    </div>
                                    <div className="d-flex align-items-center">
                                        <span className="me-2 fw-bold">Высота мачты:</span>
                                        <span className="flex-grow-1 border-bottom border-dotted"></span>
                                        <span className="ms-2">{turbine?.height} м</span>
                                    </div>
                                </div>
                            </Row>
                        </Col>
                        <Col sm className="right-column p-4">
                            <h1 style={{ textTransform: "uppercase", fontWeight: "bold", color: "#5B5B5B" }}>{turbine?.title}</h1>
                            <p>{turbine?.description}</p>
                        </Col>
                    </Row>
                </>
                : <Container style={{ display: "flex", justifyContent: "center", paddingTop: "25%" }}>
                    <BeatLoader color="#5BA1D4" />
                </Container>
        }
    </Container>
}